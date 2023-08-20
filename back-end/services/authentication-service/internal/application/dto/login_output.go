package dto

type LoginOutput struct {
	ID                          string `json:"id,omitempty"`
	Name                        string `json:"name,omitempty"`
	Email                       string `json:"email,omitempty"`
	IsMfaEnabled                bool   `json:"isMfaEnabled,omitempty"`
	PasswordVerificationTokenID string `json:"passwordVerificationToken,omitempty"`
}
