package models

type User struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	RoleId   int8   `json:"roleid"`
	Password string `json:"password"`
}
