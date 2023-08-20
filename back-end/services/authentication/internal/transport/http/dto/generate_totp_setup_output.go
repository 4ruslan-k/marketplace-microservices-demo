package dto

type TotpSetup struct {
	Image string `json:"image"`
}

type GenerateTotpSetupOutput struct {
	TotpSetup TotpSetup `json:"totpSetup"`
}
