package models

type Album struct {
	Id       int
	ArtistId int
	Name     string
	Slug     string
	Image    []byte
}

type AlbumWithArtist struct {
	Album
	Artist Artist
}
