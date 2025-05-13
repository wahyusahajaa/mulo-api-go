package dto

type CreateSongRequest struct {
	AlbumId  int    `json:"album_id" validate:"required"`
	Title    string `json:"title" validate:"required,min=1"`
	Audio    string `json:"audio" validate:"required"`
	Duration int    `json:"duration" validate:"required"`
	Image    *Image `json:"image" validate:"required"`
}

type Song struct {
	Id       int             `json:"id"`
	Title    string          `json:"title"`
	Audio    string          `json:"audio"`
	Duration int             `json:"duration"`
	Image    Image           `json:"image"`
	Album    AlbumWithArtist `json:"album"`
}
