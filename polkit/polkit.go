package polkit

import (
	"errors"
	"github.com/godbus/dbus"
	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"strconv"
)

var errAuthFailed = errors.New("authentication failed")

func NewPolKitAuthDetails(authFlags int) map[string]string {
	var details = make(map[string]string)
	details["exAuth"] = "true"
	details["exAuthFlags"] = strconv.Itoa(authFlags)
	return details
}

func CheckAuth(actionId string, busName string, details map[string]string) error {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", busName)

	ret, err := authority.CheckAuthorization(0, subject,
		actionId, details,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return err
	}

	if ret.IsAuthorized {
		return nil
	}
	return errAuthFailed
}
