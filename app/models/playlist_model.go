package models

type CreatePlaylistInput struct {
	Name   string
	UserId int
}

type Playlist struct {
	Id   int
	Name string
}
