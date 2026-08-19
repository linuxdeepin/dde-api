// SPDX-FileCopyrightText: 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/godbus/dbus/v5"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/system/org.freedesktop.dbus"
	"github.com/linuxdeepin/go-lib/dbusutil"
)

const (
	allowCallerStateFile = "/run/dde-api/locale-helper-allow-callers.json"
	privilegedGroup      = "deepin-daemon"
	invalidGroupID       = ^uint32(0)
)

type allowCallerBus interface {
	NameHasOwner(name string) (bool, error)
	GetConnUID(name string) (uint32, error)
	GetConnPID(name string) (uint32, error)
	GetBusID() (string, error)
}

type localeHelperBus struct {
	*dbusutil.Service
}

func (s localeHelperBus) GetBusID() (string, error) {
	var id string
	err := s.Conn().BusObject().Call("org.freedesktop.DBus.GetId", 0).Store(&id)
	return id, err
}

// callerInfo tracks a registered allow-caller and its UID at registration time.
// The UID is used to detect D-Bus unique name reuse attacks in authorize().
type callerInfo struct {
	uid uint32
}

// persistedCallerEntry is the per-entry JSON format for persisted allow-callers.
type persistedCallerEntry struct {
	Name string `json:"name"`
	UID  uint32 `json:"uid"`
}

type persistedAllowCallers struct {
	BusID   string                 `json:"busId"`
	Callers []persistedCallerEntry `json:"callers"`
}

// errAllowCallerNotEnabled reports that no caller has been registered through
// SetAllowCaller yet, which means the service was not started via
// deepin-security-loader. Callers should fall back to Polkit authorization.
var errAllowCallerNotEnabled = errors.New("allow-caller registry is not enabled")

type allowCallerRegistry struct {
	service           allowCallerBus
	stateFile         string
	busID             string
	privilegedGroupID uint32
	processStartTime  func(uint32) (uint64, error)
	processGroups     func(uint32) ([]uint32, error)
	processParent     func(uint32) (uint32, error)

	mu        sync.RWMutex
	callers   map[string]callerInfo
	persistMu sync.Mutex

	signalLoop *dbusutil.SignalLoop
}

func newAllowCallerRegistry(service *dbusutil.Service) (*allowCallerRegistry, error) {
	groupID, err := lookupGroupID(privilegedGroup)
	if err != nil {
		logger.Warningf("failed to resolve privileged group %s: %v", privilegedGroup, err)
		groupID = invalidGroupID
	}

	registry, err := newAllowCallerRegistryWithConfig(localeHelperBus{service}, allowCallerStateFile, groupID)
	if err != nil {
		return nil, err
	}
	if err := registry.load(); err != nil {
		logger.Warning("failed to load allow-caller state:", err)
	}

	registry.signalLoop = dbusutil.NewSignalLoop(service.Conn(), 10)
	registry.signalLoop.Start()
	dbusDaemon := ofdbus.NewDBus(service.Conn())
	dbusDaemon.InitSignalExt(registry.signalLoop, true)
	_, err = dbusDaemon.ConnectNameOwnerChanged(func(name, oldOwner, newOwner string) {
		if strings.HasPrefix(name, ":") && oldOwner != "" && newOwner == "" {
			registry.removeCaller(name)
		}
	})
	if err != nil {
		registry.signalLoop.Stop()
		registry.signalLoop = nil
		return nil, fmt.Errorf("watch allow-caller owner changes failed: %w", err)
	}

	return registry, nil
}

func newAllowCallerRegistryWithConfig(service allowCallerBus, stateFile string, privilegedGroupID uint32) (*allowCallerRegistry, error) {
	busID, err := service.GetBusID()
	if err != nil {
		return nil, fmt.Errorf("get system bus ID failed: %w", err)
	}
	if busID == "" {
		return nil, errors.New("system bus ID is empty")
	}
	return &allowCallerRegistry{
		service:           service,
		stateFile:         stateFile,
		busID:             busID,
		privilegedGroupID: privilegedGroupID,
		processStartTime:  getProcessStartTime,
		processGroups:     getProcessGroups,
		processParent:     getProcessParentPID,
		callers:           make(map[string]callerInfo),
	}, nil
}

func (r *allowCallerRegistry) close() {
	if r != nil && r.signalLoop != nil {
		r.signalLoop.Stop()
	}
}

func lookupGroupID(name string) (uint32, error) {
	group, err := user.LookupGroup(name)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(group.Gid, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid gid %q for group %s: %w", group.Gid, name, err)
	}
	return uint32(value), nil
}

// readProcStatus reads and returns the raw content of /proc/[pid]/status.
func readProcStatus(pid uint32) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// getProcessStartTime reads field 22 (starttime) from /proc/[pid]/stat.
// It is used to detect PID reuse (TOCTOU mitigation).
func getProcessStartTime(pid uint32) (uint64, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, err
	}
	// Format: pid (comm) state ppid ... starttime ...
	// Find the closing paren of comm to skip variable-length comm.
	content := string(data)
	rparen := strings.LastIndex(content, ")")
	if rparen == -1 {
		return 0, fmt.Errorf("malformed /proc/%d/stat: no closing paren", pid)
	}
	fields := strings.Fields(content[rparen+1:])
	// starttime is field 22 (0-indexed: 19 relative to after-paren start)
	// After ")" the fields are: state ppid pgrp session tty_nr tpgid flags
	// minflt cminflt majflt cmajflt utime stime cutime cstime priority nice
	// num_threads itrealvalue starttime
	if len(fields) < 20 {
		return 0, fmt.Errorf("malformed /proc/%d/stat: too few fields", pid)
	}
	// starttime is field index 19 (0-based) after the closing paren
	starttime, err := strconv.ParseUint(fields[19], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid starttime for pid %d: %w", pid, err)
	}
	return starttime, nil
}

func getProcessGroups(pid uint32) ([]uint32, error) {
	content, err := readProcStatus(pid)
	if err != nil {
		return nil, err
	}

	var groups []uint32
	for _, line := range strings.Split(content, "\n") {
		if !strings.HasPrefix(line, "Gid:") && !strings.HasPrefix(line, "Groups:") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		for _, value := range strings.Fields(parts[1]) {
			gid, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid process gid %q: %w", value, err)
			}
			groups = append(groups, uint32(gid))
		}
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("no group credentials found for pid %d", pid)
	}
	return groups, nil
}

func getProcessParentPID(pid uint32) (uint32, error) {
	content, err := readProcStatus(pid)
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(content, "\n") {
		if !strings.HasPrefix(line, "PPid:") {
			continue
		}
		fields := strings.Fields(strings.TrimPrefix(line, "PPid:"))
		if len(fields) != 1 {
			return 0, fmt.Errorf("invalid PPid entry for pid %d", pid)
		}
		parentPID, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid parent pid %q for pid %d: %w", fields[0], pid, err)
		}
		return uint32(parentPID), nil
	}
	return 0, fmt.Errorf("no PPid entry found for pid %d", pid)
}

func isProcessDescendant(pid, ancestorPID uint32, processParent func(uint32) (uint32, error)) (bool, error) {
	if pid == 0 || ancestorPID == 0 || pid == ancestorPID {
		return false, nil
	}

	visited := make(map[uint32]struct{})
	currentPID := pid
	for currentPID != 0 {
		if _, exists := visited[currentPID]; exists {
			return false, fmt.Errorf("cycle detected in process ancestry at pid %d", currentPID)
		}
		visited[currentPID] = struct{}{}

		parentPID, err := processParent(currentPID)
		if err != nil {
			return false, err
		}
		if parentPID == ancestorPID {
			return true, nil
		}
		currentPID = parentPID
	}
	return false, nil
}

func (r *allowCallerRegistry) addCaller(sender dbus.Sender, uniqueName string) error {
	if r == nil {
		return errors.New("allow-caller registry is nil")
	}
	if sender == "" {
		return errors.New("D-Bus sender is empty")
	}
	if !strings.HasPrefix(uniqueName, ":") {
		return fmt.Errorf("invalid D-Bus unique name %q", uniqueName)
	}

	hasOwner, err := r.service.NameHasOwner(uniqueName)
	if err != nil {
		return fmt.Errorf("check D-Bus owner %q failed: %w", uniqueName, err)
	}
	if !hasOwner {
		return fmt.Errorf("D-Bus caller %q has no owner", uniqueName)
	}
	if err := r.authorizeRegistrar(sender, uniqueName); err != nil {
		return err
	}

	// Capture the target's UID at registration time for later reuse detection.
	targetUID, err := r.service.GetConnUID(uniqueName)
	if err != nil {
		return fmt.Errorf("get target caller %s UID failed: %w", uniqueName, err)
	}

	r.persistMu.Lock()
	defer r.persistMu.Unlock()

	r.mu.Lock()
	if _, exists := r.callers[uniqueName]; exists {
		r.mu.Unlock()
		return nil
	}
	r.callers[uniqueName] = callerInfo{uid: targetUID}
	callers := r.snapshotCallersLocked()
	r.mu.Unlock()

	if err := r.save(callers); err != nil {
		r.mu.Lock()
		delete(r.callers, uniqueName)
		r.mu.Unlock()
		return err
	}

	logger.Infof("registered allow-caller %s for %s (uid %d)", uniqueName, dbusServiceName, targetUID)
	return nil
}

func (r *allowCallerRegistry) authorizeRegistrar(sender dbus.Sender, uniqueName string) error {
	senderUID, err := r.service.GetConnUID(string(sender))
	if err != nil {
		return fmt.Errorf("get SetAllowCaller sender %s UID failed: %w", sender, err)
	}
	if senderUID == 0 {
		return nil
	}
	if r.privilegedGroupID == invalidGroupID {
		return fmt.Errorf("privileged group %s is unavailable", privilegedGroup)
	}

	senderPID, err := r.service.GetConnPID(string(sender))
	if err != nil {
		return fmt.Errorf("get SetAllowCaller sender %s PID failed: %w", sender, err)
	}

	// Capture starttime before reading /proc to detect PID reuse (TOCTOU mitigation).
	senderStartTime, err := r.processStartTime(senderPID)
	if err != nil {
		return fmt.Errorf("get sender %s start time failed: %w", sender, err)
	}

	groups, err := r.processGroups(senderPID)
	if err != nil {
		return fmt.Errorf("get SetAllowCaller sender %s groups failed: %w", sender, err)
	}

	// Verify PID was not recycled during /proc access.
	checkStartTime, err := r.processStartTime(senderPID)
	if err != nil || checkStartTime != senderStartTime {
		return fmt.Errorf("sender %s PID %d reused during authorization", sender, senderPID)
	}

	if !containsGroup(groups, r.privilegedGroupID) {
		return fmt.Errorf("D-Bus caller %s is not in privileged group %s", sender, privilegedGroup)
	}

	targetUID, err := r.service.GetConnUID(uniqueName)
	if err != nil {
		return fmt.Errorf("get target caller %s UID failed: %w", uniqueName, err)
	}
	if targetUID != senderUID {
		return fmt.Errorf("SetAllowCaller sender UID %d does not own target %s with UID %d", senderUID, uniqueName, targetUID)
	}

	targetPID, err := r.service.GetConnPID(uniqueName)
	if err != nil {
		return fmt.Errorf("get target caller %s PID failed: %w", uniqueName, err)
	}
	isDescendant, err := isProcessDescendant(targetPID, senderPID, r.processParent)
	if err != nil {
		return fmt.Errorf("verify target caller %s process ancestry failed: %w", uniqueName, err)
	}
	if !isDescendant {
		return fmt.Errorf("target caller %s PID %d is not a descendant of SetAllowCaller sender %s PID %d",
			uniqueName, targetPID, sender, senderPID)
	}
	return nil
}

func containsGroup(groups []uint32, target uint32) bool {
	for _, group := range groups {
		if group == target {
			return true
		}
	}
	return false
}

func (r *allowCallerRegistry) authorize(sender dbus.Sender) error {
	if r == nil {
		return errors.New("allow-caller registry is nil")
	}
	if sender == "" {
		return errors.New("D-Bus sender is empty")
	}

	uid, err := r.service.GetConnUID(string(sender))
	if err != nil {
		return fmt.Errorf("get caller %s UID failed: %w", sender, err)
	}
	if uid == 0 {
		return nil
	}

	r.mu.RLock()
	info, exists := r.callers[string(sender)]
	r.mu.RUnlock()
	// Verify UID matches to prevent D-Bus unique name reuse attacks.
	if exists && info.uid == uid {
		return nil
	}
	// Caller is not in the allow-caller registry. Fall back to Polkit authorization
	// regardless of whether the registry is empty or the caller is simply not registered.
	return errAllowCallerNotEnabled
}

func (r *allowCallerRegistry) removeCaller(uniqueName string) {
	if r == nil || uniqueName == "" {
		return
	}

	r.persistMu.Lock()
	defer r.persistMu.Unlock()

	r.mu.Lock()
	if _, exists := r.callers[uniqueName]; !exists {
		r.mu.Unlock()
		return
	}
	delete(r.callers, uniqueName)
	callers := r.snapshotCallersLocked()
	r.mu.Unlock()

	if err := r.save(callers); err != nil {
		logger.Warningf("failed to persist removal of allow-caller %s: %v", uniqueName, err)
	}
}

func (r *allowCallerRegistry) snapshotCallersLocked() []persistedCallerEntry {
	entries := make([]persistedCallerEntry, 0, len(r.callers))
	for name, info := range r.callers {
		entries = append(entries, persistedCallerEntry{Name: name, UID: info.uid})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries
}

func (r *allowCallerRegistry) save(callers []persistedCallerEntry) error {
	if r.busID == "" {
		return errors.New("cannot persist allow-callers without a system bus ID")
	}
	state := persistedAllowCallers{BusID: r.busID, Callers: callers}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	dir := filepath.Dir(r.stateFile)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".locale-helper-allow-callers-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, r.stateFile)
}

func (r *allowCallerRegistry) load() error {
	r.persistMu.Lock()
	defer r.persistMu.Unlock()

	data, err := os.ReadFile(r.stateFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var state persistedAllowCallers
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	if state.BusID == "" || state.BusID != r.busID {
		return nil
	}

	loaded := make(map[string]callerInfo)
	for _, entry := range state.Callers {
		if !strings.HasPrefix(entry.Name, ":") {
			continue
		}
		hasOwner, err := r.service.NameHasOwner(entry.Name)
		if err != nil {
			logger.Warningf("failed to check persisted allow-caller %s: %v", entry.Name, err)
			continue
		}
		if hasOwner {
			loaded[entry.Name] = callerInfo{uid: entry.UID}
		}
	}

	r.mu.Lock()
	for name, info := range loaded {
		r.callers[name] = info
	}
	r.mu.Unlock()
	return nil
}
