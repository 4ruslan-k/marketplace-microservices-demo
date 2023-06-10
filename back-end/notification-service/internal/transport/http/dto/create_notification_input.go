package dto

type CreateNotificationInput struct {
	TypeID string `json:"type_id" validate:"required"`
}
