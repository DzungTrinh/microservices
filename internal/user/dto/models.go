package dto

type UserDTO struct {
	ID            string   `json:"id"`
	Email         string   `json:"email"`
	Username      string   `json:"username"`
	EmailVerified bool     `json:"email_verified"`
	Roles         []string `json:"roles"`
	Permissions   []string `json:"permissions"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}
