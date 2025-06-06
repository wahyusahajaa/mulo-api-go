package dto

type Artist struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image Image  `json:"image"`
} // @name Artist

type CreateArtistRequest struct {
	Name  string `json:"name" validate:"required"`
	Image *Image `json:"image" validate:"required"`
} // @name CreateArtistRequest
