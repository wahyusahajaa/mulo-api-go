package dto

type User struct {
	Id       int    `json:"id"`
	Fullname string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Image    Image  `json:"image"`
}

type CreateUserInput struct {
	Fullname string `json:"full_name" validate:"required"`
	Image    *Image `json:"image,omitempty"`
}
