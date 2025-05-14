package dto

type CreateGenreRequest struct {
	Name  string `json:"name" validate:"required"`
	Image *Image `json:"image" validate:"required"`
} //@name CreateGenreRequest

type Genre struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image Image  `json:"image"`
} //@name Genre
