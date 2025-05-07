package dto

type CreatePlaylistRequest struct {
	Name string `json:"name" validate:"required"`
}

type Playlist struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
