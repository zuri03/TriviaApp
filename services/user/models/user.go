package models

type User struct {
	Id        string `json:"Id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created"`
	Role      string `json:"role"`
}
