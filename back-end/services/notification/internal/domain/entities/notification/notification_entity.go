package notification

import (
	customErrors "shared/errors"
	"time"
)

var ErrEmptyTypeID = customErrors.NewIncorrectInputError("invalid_input", "Notification type id is empty")

type Notification struct {
	typeID          string
	messageTemplate string
	titleTemplate   string
	createdAt       time.Time
}

var MFAEnabledNotification = Notification{
	typeID:          "mfa-enabled-v1",
	titleTemplate:   "MFA OTP Enabled",
	messageTemplate: "Your account's Multi-Factor Authentication (MFA) One-Time Password (OTP) feature has been enabled. Your account is now more secure!",
}

var MFADisabledNotification = Notification{
	typeID:          "mfa-disabled-v1",
	titleTemplate:   "MFA OTP Disabled",
	messageTemplate: "Your account's Multi-Factor Authentication (MFA) One-Time Password (OTP) feature has been disabled. Your account is now more secure!",
}

var NotificationByTypeIds = map[string]Notification{
	MFAEnabledNotification.typeID:  MFAEnabledNotification,
	MFADisabledNotification.typeID: MFADisabledNotification,
}

func (d Notification) TypeID() string {
	return d.typeID
}

func (d Notification) MessageTemplate() string {
	return d.messageTemplate
}

func (d Notification) TitleTemplate() string {
	return d.titleTemplate
}

func (d Notification) IsZero() bool {
	return d.typeID == ""
}

func (d Notification) CreatedAt() time.Time {
	return d.createdAt
}
