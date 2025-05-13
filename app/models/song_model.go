package models

type Song struct {
	Id       int
	AlbumId  int
	Title    string
	Audio    string
	Duration int
	Image    []byte
	Album    AlbumWithArtist
}

type CreateSongInput struct {
	AlbumId  int
	Title    string
	Audio    string
	Duration int
	Image    []byte
}
