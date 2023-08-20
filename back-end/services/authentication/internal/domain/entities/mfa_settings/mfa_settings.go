package soccialaccount

import (
	"time"
)

type MfaSettings struct {
	isMfaEnabled bool
	totpSecret   string
	createdAt    time.Time
	updatedAt    time.Time
}

func NewMfaSettingsFromDatabase(
	isMfaEnabled bool,
	totpSecret string,
	createdAt time.Time,
	updatedAt time.Time,
) MfaSettings {
	mfaSettings := MfaSettings{
		isMfaEnabled: isMfaEnabled,
		totpSecret:   totpSecret,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
	return mfaSettings
}

func NewMfaSettings(isMfaEnabled bool, totpSecret string) MfaSettings {
	mfaSettings := MfaSettings{
		isMfaEnabled: isMfaEnabled,
		totpSecret:   totpSecret,
		createdAt:    time.Now(),
	}
	return mfaSettings
}

func (m MfaSettings) IsMfaEnabled() bool {
	return m.isMfaEnabled
}

func (m MfaSettings) TotpSecret() string {
	return m.totpSecret
}

func (m MfaSettings) CreatedAt() time.Time {
	return m.createdAt
}

func (m MfaSettings) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *MfaSettings) SetTotpSecret(totpSecret string) {
	m.totpSecret = totpSecret
}

func (m *MfaSettings) SetMfaStatus(isMfaEnabled bool) {
	m.isMfaEnabled = isMfaEnabled
}
