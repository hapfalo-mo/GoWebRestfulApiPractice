package dto

type UserResponse struct {
	UserID      int     `json:"ID"`
	UserName    string  `json:"UserName"`
	Email       string  `json:"Email"`
	PhoneNumber string  `json:"PhoneNumber"`
	CreatedAt   string  `json:"CreatedAt"`
	UpdatedAt   string  `json:"UpdatedAt"`
	DeletedAt   *string `json:"DeletedAt"`
	OrderCount  int     `json:"OrderCount"`
}
