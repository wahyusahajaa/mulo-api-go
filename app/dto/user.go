package dto

type User struct {
	Id       int    `json:"id"`
	Fullname string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    Image  `json:"image"`
}

type UserUpdateInput struct {
	Fullname string `json:"full_name"`
	Image    *Image `json:"image,omitempty"`
}
