package models

type User struct {
	UserID          int     `json:"user_id"`
	Username        string  `json:"username"`
	UserEmail       string  `json:"user_email"`
	UserPhoneNumber string  `json:"user_phone_number"`
	UserCreatedAt   string  `json:"user_created_at"`
	UserUpdatedAt   string  `json:"user_updated_at"`
	UserDeletedAt   *string `json:"user_deleted_at"`
}
