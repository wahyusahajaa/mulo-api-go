package dto

type CreateAlbumRequest struct {
	Name     string `json:"name" validate:"required"`
	ArtistId int    `json:"artist_id" validate:"required"`
	Image    *Image `json:"image,omitempty"`
}

type Album struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	Image  Image  `json:"image"`
	Artist Artist `json:"artist"`
}
