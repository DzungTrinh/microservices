package events

// UserRegisteredPayload represents the payload for UserRegistered event.
type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
}

// AdminUserCreatedPayload represents the payload for AdminUserCreated event.
type AdminUserCreatedPayload struct {
	UserID string `json:"user_id"`
}
