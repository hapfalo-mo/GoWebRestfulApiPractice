package dtorequest

type LoginRequest struct {
	UserEmail string `json:"user_email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
