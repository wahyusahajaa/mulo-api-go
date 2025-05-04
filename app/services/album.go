package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type albumService struct {
	repo contracts.AlbumRepository
	log  *logrus.Logger
}

func NewAlbumService(repo contracts.AlbumRepository, log *logrus.Logger) contracts.AlbumService {
	return &albumService{
		repo: repo,
		log:  log,
	}
}

func (svc *albumService) GetAll(ctx context.Context, pageSize int, offset int) (albums []models.Album, err error) {
	albums, err = svc.repo.FindAll(ctx, pageSize, offset)

	if err != nil {
		svc.log.WithError(err).Error("error in album service")
		return
	}

	return
}

func (svc *albumService) GetAlbumById(ctx context.Context, id int) (album *models.AlbumWithArtist, err error) {
	album, err = svc.repo.FindAlbumById(ctx, id)

	if err != nil {
		svc.log.WithError(err).Error("error in album service")
		return
	}

	return
}

func (svc *albumService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.Count(ctx)

	if err != nil {
		svc.log.WithError(err).Error("error in album service")
		return
	}

	return
}

func (svc *albumService) Create(ctx context.Context, artistId int, name string, slug string, image []byte) (err error) {
	err = svc.repo.Store(ctx, artistId, name, slug, image)

	if err != nil {
		svc.log.WithError(err).Error("error in album service")
		return
	}

	return
}

func (svc *albumService) CheckDuplicateAlbumBySlug(ctx context.Context, slug string) (exists bool, err error) {
	exists, err = svc.repo.FindDuplicateAlbumBySlug(ctx, slug)

	if err != nil {
		svc.log.WithError(err).Error("error in album service")
		return
	}

	return
}

func (svc *albumService) UpdateAlbum(ctx context.Context, artistId int, name string, slug string, image []byte, id int) (err error) {
	err = svc.repo.Update(ctx, artistId, name, slug, image, id)

	if err != nil {
		svc.log.WithError(err).Error("error in album service")
		return
	}

	return
}

func (svc *albumService) DeleteAlbum(ctx context.Context, id int) (err error) {
	err = svc.repo.Delete(ctx, id)

	if err != nil {
		svc.log.WithError(err).Error("error iin album service")
		return
	}

	return
}
