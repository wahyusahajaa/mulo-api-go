package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/mocks"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

type fields struct {
	repoCountResult int
	repoCountErr    error
	repoListResult  []models.Genre
	repoListErr     error
}

type expected struct {
	result []dto.Genre
	total  int
	err    error
}

type tests struct {
	name   string
	fields fields
	want   expected
}

func TestGenreService_GetAll(t *testing.T) {
	// Make json byte from dto.image
	image1 := dto.Image{Src: "image1.png", BlurHash: "blur1"}
	image2 := dto.Image{Src: "image2.png", BlurHash: "blur2"}
	// Parse image to bytes
	image1Bytes, _ := json.Marshal(image1)
	image2Bytes, _ := json.Marshal(image2)

	testsLists := []tests{
		{
			name: "success",
			fields: fields{
				repoCountResult: 2,
				repoCountErr:    nil,
				repoListResult: []models.Genre{
					{Id: 1, Name: "Genre 1", Image: image1Bytes},
					{Id: 2, Name: "Genre 2", Image: image2Bytes},
				},
				repoListErr: nil,
			},
			want: expected{
				result: []dto.Genre{
					{Id: 1, Name: "Genre 1", Image: utils.ParseImageToJSON(image1Bytes)},
					{Id: 2, Name: "Genre 2", Image: utils.ParseImageToJSON(image2Bytes)},
				},
				total: 2,
				err:   nil,
			},
		},
		{
			name: "findCountError",
			fields: fields{
				repoCountResult: 0,
				repoCountErr:    errors.New("db connection failure"),
			},
			want: expected{
				total:  0,
				result: nil,
				err:    errors.New("db connection failure"),
			},
		},
		{
			name: "findAllError",
			fields: fields{
				repoCountResult: 10,
				repoCountErr:    nil,
				repoListResult:  nil,
				repoListErr:     errors.New("db list error"),
			},
			want: expected{
				total:  0,
				result: nil,
				err:    errors.New("db list error"),
			},
		},
	}

	for _, tl := range testsLists {

		t.Run(tl.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockGenreRepository)
			svc := NewGenreService(mockRepo, nil, nil, nil)

			ctx := context.Background()
			pageSize := 10
			offset := 0

			mockRepo.On("FindCount", ctx).Return(tl.fields.repoCountResult, tl.fields.repoCountErr)

			if tl.fields.repoCountErr == nil {
				mockRepo.On("FindAll", ctx, pageSize, offset).Return(tl.fields.repoListResult, tl.fields.repoListErr)
			}

			// Actual
			result, total, err := svc.GetAll(ctx, pageSize, offset)

			// Assert
			if tl.want.err != nil {
				assert.EqualError(t, err, tl.want.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tl.want.total, total)
			assert.Equal(t, tl.want.result, result)

			mockRepo.AssertExpectations(t)
		})

	}
}
