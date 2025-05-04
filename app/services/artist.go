package services

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wahyusahajaa/mulo-api-go/app/contracts"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
)

type artistService struct {
	repo contracts.ArtistRepository
	log  *logrus.Logger
}

func NewArtistService(repo contracts.ArtistRepository, log *logrus.Logger) contracts.ArtistService {
	return &artistService{
		repo: repo,
		log:  log,
	}
}

func (svc *artistService) GetAll(ctx context.Context, pageSize, offset int) (artists []models.Artist, err error) {
	artists, err = svc.repo.FindAll(ctx, pageSize, offset)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return nil, err
	}

	return artists, nil
}

func (svc *artistService) GetArtistByIds(ctx context.Context, inClause string, artistIds []any) (artists []models.Artist, err error) {
	artists, err = svc.repo.FindByArtistIds(ctx, inClause, artistIds)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) GetCount(ctx context.Context) (total int, err error) {
	total, err = svc.repo.Count(ctx)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) Create(ctx context.Context, name string, slug string, image []byte) (err error) {
	err = svc.repo.Store(ctx, name, slug, image)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) CheckDuplicateArtistBySlug(ctx context.Context, slug string) (exists bool, err error) {
	exists, err = svc.repo.FindDuplicateArtistBySlug(ctx, slug)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) GetArtistById(ctx context.Context, artistId int) (artist *models.Artist, err error) {
	artist, err = svc.repo.FindArtistById(ctx, artistId)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) Update(ctx context.Context, name string, slug string, image []byte, artistId int) (err error) {
	err = svc.repo.Update(ctx, name, slug, image, artistId)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}

func (svc *artistService) Delete(ctx context.Context, artistId int) (err error) {
	err = svc.repo.Delete(ctx, artistId)

	if err != nil {
		svc.log.WithError(err).Error("error in artist service")
		return
	}

	return
}
