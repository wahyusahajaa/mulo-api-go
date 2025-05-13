package models

type Artist struct {
	Id    int
	Name  string
	Slug  string
	Image []byte
}

type CreateArtistInput struct {
	Name  string
	Slug  string
	Image []byte
}
