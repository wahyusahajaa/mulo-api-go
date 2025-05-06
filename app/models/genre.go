package models

type Genre struct {
	Id    int
	Name  string
	Image []byte
}

type CreateGenreInput struct {
	Name  string
	Image []byte
}
