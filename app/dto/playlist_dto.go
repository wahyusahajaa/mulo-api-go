package dto

type CreatePlaylistRequest struct {
	Name string `json:"name" validate:"required"`
} // @name CreatePlaylistRequest

type Playlist struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
} // @name Playlist

type PlaylistWithSongs struct {
	Playlist
	Song []Song `json:"songs"`
} // @name PlaylistWithSongs
