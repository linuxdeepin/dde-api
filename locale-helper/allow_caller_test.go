// SPDX-FileCopyrightText: 2026 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/godbus/dbus/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockAllowCallerBus implements allowCallerBus for testing without a real D-Bus.
type mockAllowCallerBus struct {
	busID        string
	nameHasOwner map[string]bool
	connUID      map[string]uint32
	connPID      map[string]uint32
}

func newMockBus() *mockAllowCallerBus {
	return &mockAllowCallerBus{
		busID:        "mock-bus-id",
		nameHasOwner: make(map[string]bool),
		connUID:      make(map[string]uint32),
		connPID:      make(map[string]uint32),
	}
}

func (m *mockAllowCallerBus) NameHasOwner(name string) (bool, error) {
	has, ok := m.nameHasOwner[name]
	if !ok {
		return false, fmt.Errorf("unknown name %q", name)
	}
	return has, nil
}

func (m *mockAllowCallerBus) GetConnUID(name string) (uint32, error) {
	uid, ok := m.connUID[name]
	if !ok {
		return 0, fmt.Errorf("unknown name %q", name)
	}
	return uid, nil
}

func (m *mockAllowCallerBus) GetConnPID(name string) (uint32, error) {
	pid, ok := m.connPID[name]
	if !ok {
		return 0, fmt.Errorf("unknown name %q", name)
	}
	return pid, nil
}

func (m *mockAllowCallerBus) GetBusID() (string, error) {
	return m.busID, nil
}

// newRegistryForTest creates an allowCallerRegistry backed by a mock bus and
// a temporary state file.
func newRegistryForTest(t *testing.T, bus *mockAllowCallerBus) *allowCallerRegistry {
	t.Helper()
	stateFile := filepath.Join(t.TempDir(), "allow-callers.json")
	r, err := newAllowCallerRegistryWithConfig(bus, stateFile, invalidGroupID)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, bus.busID, r.busID)
	return r
}

// ---------------------------------------------------------------------------
// containsGroup
// ---------------------------------------------------------------------------

func TestContainsGroup(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		assert.True(t, containsGroup([]uint32{1, 42, 100}, 42))
	})
	t.Run("not contains", func(t *testing.T) {
		assert.False(t, containsGroup([]uint32{1, 2, 3}, 42))
	})
	t.Run("empty slice", func(t *testing.T) {
		assert.False(t, containsGroup(nil, 42))
	})
}

// ---------------------------------------------------------------------------
// isProcessDescendant
// ---------------------------------------------------------------------------

func TestIsProcessDescendant(t *testing.T) {
	t.Run("self is not descendant", func(t *testing.T) {
		result, err := isProcessDescendant(5, 5, nil)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("pid 0 is not descendant", func(t *testing.T) {
		result, err := isProcessDescendant(0, 5, nil)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("ancestor 0 is not descendant", func(t *testing.T) {
		result, err := isProcessDescendant(5, 0, nil)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("direct child", func(t *testing.T) {
		parent := func(pid uint32) (uint32, error) {
			switch pid {
			case 10:
				return 5, nil
			}
			return 0, nil
		}
		result, err := isProcessDescendant(10, 5, parent)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("indirect descendant", func(t *testing.T) {
		parent := func(pid uint32) (uint32, error) {
			switch pid {
			case 10:
				return 7, nil
			case 7:
				return 5, nil
			}
			return 0, nil
		}
		result, err := isProcessDescendant(10, 5, parent)
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("not descendant", func(t *testing.T) {
		parent := func(pid uint32) (uint32, error) {
			switch pid {
			case 10:
				return 3, nil
			case 3:
				return 1, nil
			}
			return 0, nil
		}
		result, err := isProcessDescendant(10, 5, parent)
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("cycle detection", func(t *testing.T) {
		parent := func(pid uint32) (uint32, error) {
			switch pid {
			case 10:
				return 7, nil
			case 7:
				return 10, nil
			}
			return 0, nil
		}
		_, err := isProcessDescendant(10, 5, parent)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cycle")
	})

	t.Run("parent lookup error", func(t *testing.T) {
		parent := func(pid uint32) (uint32, error) {
			return 0, fmt.Errorf("mock error")
		}
		_, err := isProcessDescendant(10, 5, parent)
		assert.Error(t, err)
	})
}

// ---------------------------------------------------------------------------
// snapshotCallersLocked
// ---------------------------------------------------------------------------

func TestSnapshotCallersLocked(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	r.callers[":1.100"] = callerInfo{uid: 1000}
	r.callers[":1.50"] = callerInfo{uid: 1000}
	r.callers[":1.200"] = callerInfo{uid: 1000}

	result := r.snapshotCallersLocked()
	assert.Equal(t, ":1.100", result[0].Name)
	assert.Equal(t, ":1.200", result[1].Name)
	assert.Equal(t, ":1.50", result[2].Name)
	assert.Equal(t, uint32(1000), result[0].UID)
	assert.Equal(t, uint32(1000), result[1].UID)
	assert.Equal(t, uint32(1000), result[2].UID)
}

// ---------------------------------------------------------------------------
// newAllowCallerRegistryWithConfig
// ---------------------------------------------------------------------------

func TestNewAllowCallerRegistryWithConfig(t *testing.T) {
	t.Run("empty bus ID returns error", func(t *testing.T) {
		bus := newMockBus()
		bus.busID = ""
		r, err := newAllowCallerRegistryWithConfig(bus, "/tmp/state.json", 42)
		assert.Error(t, err)
		assert.Nil(t, r)
	})

	t.Run("success", func(t *testing.T) {
		bus := newMockBus()
		r, err := newAllowCallerRegistryWithConfig(bus, "/tmp/state.json", 42)
		assert.NoError(t, err)
		assert.NotNil(t, r)
		assert.Equal(t, "mock-bus-id", r.busID)
		assert.Equal(t, uint32(42), r.privilegedGroupID)
		assert.NotNil(t, r.callers)
		assert.Len(t, r.callers, 0)
		assert.NotNil(t, r.processGroups)
		assert.NotNil(t, r.processParent)
	})
}

// ---------------------------------------------------------------------------
// addCaller
// ---------------------------------------------------------------------------

func TestAddCallerNilRegistry(t *testing.T) {
	var r *allowCallerRegistry
	err := r.addCaller(dbus.Sender(":1.0"), ":1.100")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestAddCallerEmptySender(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	err := r.addCaller(dbus.Sender(""), ":1.100")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestAddCallerInvalidUniqueName(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	err := r.addCaller(dbus.Sender(":1.0"), "not-a-unique-name")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid D-Bus unique name")
}

func TestAddCallerNoOwner(t *testing.T) {
	bus := newMockBus()
	bus.nameHasOwner[":1.100"] = false
	r := newRegistryForTest(t, bus)
	err := r.addCaller(dbus.Sender(":1.0"), ":1.100")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no owner")
}

func TestAddCallerRootBypass(t *testing.T) {
	bus := newMockBus()
	bus.nameHasOwner[":1.100"] = true
	bus.connUID[":1.0"] = 0
	bus.connUID[":1.100"] = 1000
	r := newRegistryForTest(t, bus)

	err := r.addCaller(dbus.Sender(":1.0"), ":1.100")
	assert.NoError(t, err)

	r.mu.RLock()
	info, exists := r.callers[":1.100"]
	r.mu.RUnlock()
	assert.True(t, exists)
	assert.Equal(t, uint32(1000), info.uid)
}

func TestAddCallerDuplicate(t *testing.T) {
	bus := newMockBus()
	bus.nameHasOwner[":1.100"] = true
	bus.connUID[":1.0"] = 0
	bus.connUID[":1.100"] = 1000
	r := newRegistryForTest(t, bus)

	err := r.addCaller(dbus.Sender(":1.0"), ":1.100")
	require.NoError(t, err)

	err = r.addCaller(dbus.Sender(":1.0"), ":1.100")
	assert.NoError(t, err)

	r.mu.RLock()
	assert.Len(t, r.callers, 1)
	r.mu.RUnlock()
}

func TestAddCallerPersistenceRollback(t *testing.T) {
	bus := newMockBus()
	bus.nameHasOwner[":1.100"] = true
	bus.connUID[":1.0"] = 0
	bus.connUID[":1.100"] = 1000
	r := newRegistryForTest(t, bus)
	r.busID = "" // trigger save failure

	err := r.addCaller(dbus.Sender(":1.0"), ":1.100")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot persist")

	r.mu.RLock()
	_, exists := r.callers[":1.100"]
	r.mu.RUnlock()
	assert.False(t, exists, "caller must be removed from map after save failure")
}

func TestAddCallerAuthorizedSender(t *testing.T) {
	bus := newMockBus()
	senderName := dbus.Sender(":1.0")
	targetName := ":1.100"
	bus.nameHasOwner[targetName] = true
	bus.connUID[string(senderName)] = 1000
	bus.connUID[targetName] = 1000
	bus.connPID[string(senderName)] = 5
	bus.connPID[targetName] = 10

	r := newRegistryForTest(t, bus)
	r.privilegedGroupID = 42
	r.processGroups = func(pid uint32) ([]uint32, error) {
		return []uint32{42}, nil
	}
	r.processParent = func(pid uint32) (uint32, error) {
		switch pid {
		case 10:
			return 5, nil
		case 5:
			return 1, nil
		}
		return 0, nil
	}

	err := r.addCaller(senderName, targetName)
	assert.NoError(t, err)

	r.mu.RLock()
	info, exists := r.callers[targetName]
	r.mu.RUnlock()
	assert.True(t, exists)
	assert.Equal(t, uint32(1000), info.uid)
}

func TestAddCallerUnauthorizedSender(t *testing.T) {
	bus := newMockBus()
	senderName := dbus.Sender(":1.0")
	targetName := ":1.100"
	bus.nameHasOwner[targetName] = true
	bus.connUID[string(senderName)] = 1000

	r := newRegistryForTest(t, bus)
	// privilegedGroupID is invalidGroupID → non-root sender fails

	err := r.addCaller(senderName, targetName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "privileged group")
}

// ---------------------------------------------------------------------------
// authorize
// ---------------------------------------------------------------------------

func TestAuthorizeNilRegistry(t *testing.T) {
	var r *allowCallerRegistry
	err := r.authorize(dbus.Sender(":1.0"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestAuthorizeEmptySender(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	err := r.authorize(dbus.Sender(""))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestAuthorizeRootBypass(t *testing.T) {
	bus := newMockBus()
	bus.connUID[":1.0"] = 0
	r := newRegistryForTest(t, bus)

	err := r.authorize(dbus.Sender(":1.0"))
	assert.NoError(t, err)
}

func TestAuthorizeAllowedCaller(t *testing.T) {
	bus := newMockBus()
	bus.connUID[":1.0"] = 1000
	r := newRegistryForTest(t, bus)
	r.mu.Lock()
	r.callers[":1.0"] = callerInfo{uid: 1000}
	r.mu.Unlock()

	err := r.authorize(dbus.Sender(":1.0"))
	assert.NoError(t, err)
}

func TestAuthorizeUIDMismatch(t *testing.T) {
	bus := newMockBus()
	bus.connUID[":1.0"] = 2000
	r := newRegistryForTest(t, bus)
	r.mu.Lock()
	r.callers[":1.0"] = callerInfo{uid: 1000} // registered with different UID
	r.mu.Unlock()

	err := r.authorize(dbus.Sender(":1.0"))
	assert.ErrorIs(t, err, errAllowCallerNotEnabled)
}

func TestAuthorizeUnauthorizedCaller(t *testing.T) {
	bus := newMockBus()
	bus.connUID[":1.0"] = 1000
	r := newRegistryForTest(t, bus)

	// Empty registry: not launched via deepin-security-loader, fall back to Polkit.
	err := r.authorize(dbus.Sender(":1.0"))
	assert.ErrorIs(t, err, errAllowCallerNotEnabled)

	// Registry has other callers, this one is not registered either — also fall back.
	r.mu.Lock()
	r.callers[":1.9"] = callerInfo{uid: 1000}
	r.mu.Unlock()
	err = r.authorize(dbus.Sender(":1.0"))
	assert.ErrorIs(t, err, errAllowCallerNotEnabled)
}

// ---------------------------------------------------------------------------
// removeCaller
// ---------------------------------------------------------------------------

func TestRemoveCallerNilRegistry(t *testing.T) {
	var r *allowCallerRegistry
	r.removeCaller(":1.100") // must not panic
}

func TestRemoveCallerEmptyName(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	r.removeCaller("") // must not panic
}

func TestRemoveCallerNonExistent(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	r.mu.Lock()
	r.callers[":1.50"] = callerInfo{uid: 1000}
	r.mu.Unlock()

	r.removeCaller(":1.100")

	r.mu.RLock()
	assert.Len(t, r.callers, 1)
	_, exists := r.callers[":1.50"]
	r.mu.RUnlock()
	assert.True(t, exists)
}

func TestRemoveCallerExisting(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	r.mu.Lock()
	r.callers[":1.100"] = callerInfo{uid: 1000}
	r.mu.Unlock()

	r.removeCaller(":1.100")

	r.mu.RLock()
	_, exists := r.callers[":1.100"]
	r.mu.RUnlock()
	assert.False(t, exists)
}

// ---------------------------------------------------------------------------
// save / load
// ---------------------------------------------------------------------------

func TestSaveEmptyBusID(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	r.busID = ""
	err := r.save([]persistedCallerEntry{{Name: ":1.100", UID: 1000}})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot persist")
}

func TestSaveAndLoad(t *testing.T) {
	stateFile := filepath.Join(t.TempDir(), "state.json")

	// Save callers from registry r1.
	bus1 := newMockBus()
	r1, err := newAllowCallerRegistryWithConfig(bus1, stateFile, invalidGroupID)
	require.NoError(t, err)

	err = r1.save([]persistedCallerEntry{
		{Name: ":1.100", UID: 1000},
		{Name: ":1.200", UID: 2000},
	})
	assert.NoError(t, err)

	t.Run("different bus ID skips stale state", func(t *testing.T) {
		bus2 := newMockBus()
		bus2.busID = "different-bus"
		bus2.nameHasOwner[":1.100"] = true
		bus2.nameHasOwner[":1.200"] = true
		r2, err := newAllowCallerRegistryWithConfig(bus2, stateFile, invalidGroupID)
		require.NoError(t, err)

		err = r2.load()
		assert.NoError(t, err)

		r2.mu.RLock()
		assert.Len(t, r2.callers, 0)
		r2.mu.RUnlock()
	})

	t.Run("same bus ID restores callers", func(t *testing.T) {
		bus3 := newMockBus()
		bus3.nameHasOwner[":1.100"] = true
		bus3.nameHasOwner[":1.200"] = true
		r3, err := newAllowCallerRegistryWithConfig(bus3, stateFile, invalidGroupID)
		require.NoError(t, err)

		err = r3.load()
		assert.NoError(t, err)

		r3.mu.RLock()
		assert.Len(t, r3.callers, 2)
		info100, has100 := r3.callers[":1.100"]
		info200, has200 := r3.callers[":1.200"]
		r3.mu.RUnlock()
		assert.True(t, has100)
		assert.Equal(t, uint32(1000), info100.uid)
		assert.True(t, has200)
		assert.Equal(t, uint32(2000), info200.uid)
	})
}

func TestLoadMissingFile(t *testing.T) {
	r := newRegistryForTest(t, newMockBus())
	err := r.load()
	assert.NoError(t, err)
	assert.Len(t, r.callers, 0)
}

func TestLoadStaleCallerSkipped(t *testing.T) {
	stateFile := filepath.Join(t.TempDir(), "state.json")

	// Persist :1.100 and :1.200, then :1.200 disappears.
	bus := newMockBus()
	r, err := newAllowCallerRegistryWithConfig(bus, stateFile, invalidGroupID)
	require.NoError(t, err)
	err = r.save([]persistedCallerEntry{
		{Name: ":1.100", UID: 1000},
		{Name: ":1.200", UID: 2000},
	})
	require.NoError(t, err)

	bus2 := newMockBus()
	bus2.busID = bus.busID
	bus2.nameHasOwner[":1.100"] = true
	bus2.nameHasOwner[":1.200"] = false
	r2, err := newAllowCallerRegistryWithConfig(bus2, stateFile, invalidGroupID)
	require.NoError(t, err)

	err = r2.load()
	assert.NoError(t, err)

	r2.mu.RLock()
	assert.Len(t, r2.callers, 1)
	info100, has100 := r2.callers[":1.100"]
	_, has200 := r2.callers[":1.200"]
	r2.mu.RUnlock()
	assert.True(t, has100)
	assert.Equal(t, uint32(1000), info100.uid)
	assert.False(t, has200, "stale caller without owner must be filtered out")
}

// ---------------------------------------------------------------------------
// authorizeRegistrar (direct coverage of the authorization gate)
// ---------------------------------------------------------------------------

func TestAuthorizeRegistrar(t *testing.T) {
	t.Run("root sender bypass", func(t *testing.T) {
		bus := newMockBus()
		bus.connUID[":1.0"] = 0
		r := newRegistryForTest(t, bus)
		err := r.authorizeRegistrar(dbus.Sender(":1.0"), ":1.100")
		assert.NoError(t, err)
	})

	t.Run("unavailable privileged group", func(t *testing.T) {
		bus := newMockBus()
		bus.connUID[":1.0"] = 1000
		r := newRegistryForTest(t, bus)
		// r.privilegedGroupID is invalidGroupID
		err := r.authorizeRegistrar(dbus.Sender(":1.0"), ":1.100")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "privileged group")
	})

	t.Run("sender not in privileged group", func(t *testing.T) {
		bus := newMockBus()
		bus.connUID[":1.0"] = 1000
		bus.connPID[":1.0"] = 5
		r := newRegistryForTest(t, bus)
		r.privilegedGroupID = 42
		r.processGroups = func(pid uint32) ([]uint32, error) {
			return []uint32{99}, nil
		}
		err := r.authorizeRegistrar(dbus.Sender(":1.0"), ":1.100")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not in privileged group")
	})

	t.Run("target UID mismatch", func(t *testing.T) {
		bus := newMockBus()
		bus.connUID[":1.0"] = 1000
		bus.connPID[":1.0"] = 5
		bus.connUID[":1.100"] = 2000
		r := newRegistryForTest(t, bus)
		r.privilegedGroupID = 42
		r.processGroups = func(pid uint32) ([]uint32, error) {
			return []uint32{42}, nil
		}
		err := r.authorizeRegistrar(dbus.Sender(":1.0"), ":1.100")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not own target")
	})

	t.Run("target not a descendant", func(t *testing.T) {
		bus := newMockBus()
		bus.connUID[":1.0"] = 1000
		bus.connPID[":1.0"] = 5
		bus.connUID[":1.100"] = 1000
		bus.connPID[":1.100"] = 10
		r := newRegistryForTest(t, bus)
		r.privilegedGroupID = 42
		r.processGroups = func(pid uint32) ([]uint32, error) {
			return []uint32{42}, nil
		}
		r.processParent = func(pid uint32) (uint32, error) {
			switch pid {
			case 10:
				return 3, nil
			case 3:
				return 1, nil
			case 5:
				return 1, nil
			}
			return 0, nil
		}
		err := r.authorizeRegistrar(dbus.Sender(":1.0"), ":1.100")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a descendant")
	})

	t.Run("successful authorization", func(t *testing.T) {
		bus := newMockBus()
		bus.connUID[":1.0"] = 1000
		bus.connPID[":1.0"] = 5
		bus.connUID[":1.100"] = 1000
		bus.connPID[":1.100"] = 10
		r := newRegistryForTest(t, bus)
		r.privilegedGroupID = 42
		r.processGroups = func(pid uint32) ([]uint32, error) {
			return []uint32{42}, nil
		}
		r.processParent = func(pid uint32) (uint32, error) {
			switch pid {
			case 10:
				return 5, nil
			case 5:
				return 1, nil
			}
			return 0, nil
		}
		err := r.authorizeRegistrar(dbus.Sender(":1.0"), ":1.100")
		assert.NoError(t, err)
	})
}

// ---------------------------------------------------------------------------
// lookupGroupID (OS integration smoke test)
// ---------------------------------------------------------------------------

func TestLookupGroupID(t *testing.T) {
	_, err := lookupGroupID("root")
	if err != nil {
		t.Skipf("group 'root' not found on this system: %v; skipping", err)
	}
	// If the root group exists, verify it has GID 0.
	gid, err := lookupGroupID("root")
	require.NoError(t, err)
	assert.Equal(t, uint32(0), gid)
}