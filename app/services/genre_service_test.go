package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wahyusahajaa/mulo-api-go/app/dto"
	"github.com/wahyusahajaa/mulo-api-go/app/mocks"
	"github.com/wahyusahajaa/mulo-api-go/app/models"
	"github.com/wahyusahajaa/mulo-api-go/pkg/errs"
	"github.com/wahyusahajaa/mulo-api-go/pkg/utils"
)

func TestGenreService_GetAll(t *testing.T) {
	// Make json byte from dto.image
	image1 := dto.Image{Src: "image1.png", BlurHash: "blur1"}
	image2 := dto.Image{Src: "image2.png", BlurHash: "blur2"}
	// Parse image to bytes
	image1Bytes, _ := json.Marshal(image1)
	image2Bytes, _ := json.Marshal(image2)

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

	type scenario struct {
		name     string
		fields   fields
		expected expected
	}

	testCases := []scenario{
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
			expected: expected{
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
			expected: expected{
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
			expected: expected{
				total:  0,
				result: nil,
				err:    errors.New("db list error"),
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(mocks.MockGenreRepository)
			svc := NewGenreService(mockRepo, nil, nil, nil)

			ctx := context.Background()
			pageSize := 10
			offset := 0

			mockRepo.On("FindCount", ctx).Return(tc.fields.repoCountResult, tc.fields.repoCountErr)

			if tc.fields.repoCountErr == nil {
				mockRepo.On("FindAll", ctx, pageSize, offset).Return(tc.fields.repoListResult, tc.fields.repoListErr)
			}

			// Actual
			result, total, err := svc.GetAll(ctx, pageSize, offset)

			// Assert
			if tc.expected.err != nil {
				assert.Equal(t, tc.expected.err, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected.total, total)
			assert.Equal(t, tc.expected.result, result)

			mockRepo.AssertExpectations(t)
		})

	}
}

func TestGenreService_GetGenreById(t *testing.T) {
	type fields struct {
		repoErr    error
		repoResult *models.Genre
		genreId    int
	}

	type expected struct {
		result dto.Genre
		err    error
	}

	type scenario struct {
		name     string
		fields   fields
		expected expected
	}

	image1 := dto.Image{Src: "image1.png", BlurHash: "image1"}
	image1_bytes, _ := json.Marshal(image1)

	testCases := []scenario{
		{
			name: "success",
			fields: fields{
				repoErr:    nil,
				repoResult: &models.Genre{Id: 1, Name: "Genre 1", Image: image1_bytes},
				genreId:    1,
			},
			expected: expected{
				result: dto.Genre{Id: 1, Name: "Genre 1", Image: image1},
				err:    nil,
			},
		},
		{
			name: "genreNotFound",
			fields: fields{
				repoErr:    errs.NewNotFoundError("Genre", "id", 1),
				repoResult: nil,
				genreId:    1,
			},
			expected: expected{
				result: dto.Genre{},
				err:    errs.NewNotFoundError("Genre", "id", 1),
			},
		},
		{
			name: "error",
			fields: fields{
				repoErr:    errors.New("database failure"),
				repoResult: nil,
				genreId:    1,
			},
			expected: expected{
				result: dto.Genre{},
				err:    errors.New("database failure"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockGenreRepository)
			svc := NewGenreService(mockRepo, nil, nil, nil)
			ctx := context.Background()

			mockRepo.On("FindGenreById", mock.Anything, tc.fields.genreId).Return(tc.fields.repoResult, tc.fields.repoErr)
			actualResult, err := svc.GetGenreById(ctx, tc.fields.genreId)

			if tc.expected.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expected.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected.result, actualResult)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenreService_CreateGenre(t *testing.T) {
	image1 := dto.Image{Src: "image1.png", BlurHash: "image1"}
	imageBytes, _ := utils.ParseImageToByte(&image1)
	var errValidation = errors.New("validation failed")

	type fields struct {
		repoStoreErr error
		reqDtoGenre  dto.CreateGenreRequest
	}

	type expected struct {
		errType      error
		valErrorMaps map[string]string
	}

	type scenario struct {
		name     string
		fields   fields
		expected expected
	}

	testCases := []scenario{
		{
			name: "success",
			fields: fields{
				repoStoreErr: nil,
				reqDtoGenre:  dto.CreateGenreRequest{Name: "Genre 1", Image: &image1},
			},
			expected: expected{
				errType:      nil,
				valErrorMaps: nil,
			},
		},
		{
			name: "repoStoreErr",
			fields: fields{
				repoStoreErr: errors.New("database failure"),
				reqDtoGenre:  dto.CreateGenreRequest{Name: "Genre 1", Image: &image1},
			},
			expected: expected{
				errType:      errors.New("database failure"),
				valErrorMaps: nil,
			},
		},
		{
			name: "validationErrors_RequiredName",
			fields: fields{
				reqDtoGenre: dto.CreateGenreRequest{Name: "", Image: &image1},
			},
			expected: expected{
				errType:      errValidation,
				valErrorMaps: map[string]string{"name": "Field is required"},
			},
		},
		{
			name: "validationErrors_RequiredImage",
			fields: fields{
				reqDtoGenre: dto.CreateGenreRequest{Name: "Genre 1", Image: nil},
			},
			expected: expected{
				errType:      errValidation,
				valErrorMaps: map[string]string{"image": "Field is required"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockGenreRepository)
			svc := NewGenreService(mockRepo, nil, nil, nil)
			ctx := context.Background()

			if tc.expected.valErrorMaps == nil {
				// convert DTO to input model
				input := models.CreateGenreInput{
					Name:  tc.fields.reqDtoGenre.Name,
					Image: imageBytes,
				}
				mockRepo.On("Store", mock.Anything, input).Return(tc.fields.repoStoreErr)
			}

			err := svc.CreateGenre(ctx, tc.fields.reqDtoGenre)

			if tc.expected.errType == nil {
				assert.NoError(t, err)
			} else if errors.Is(tc.expected.errType, errValidation) {
				var valErr *errs.BadRequestError
				assert.ErrorAs(t, err, &valErr)
				assert.Equal(t, tc.expected.valErrorMaps, valErr.Errors)
			} else {
				assert.Equal(t, tc.expected.errType, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}

}

func TestGenreService_UpdateGenre(t *testing.T) {
	image1 := dto.Image{Src: "image1.png", BlurHash: "abs"}
	image1Bytes, _ := json.Marshal(image1)
	var errValidation = errors.New("validation failed")

	type fields struct {
		repoGenreErr    error
		repoGenreResult bool
		repoUpdateErr   error
		reqDtoGenre     dto.CreateGenreRequest
	}

	type expected struct {
		errType      error
		valErrorMaps map[string]string
	}

	type scenario struct {
		name     string
		fields   fields
		expected expected
	}

	testCases := []scenario{
		{
			name: "success",
			fields: fields{
				repoGenreErr:    nil,
				repoUpdateErr:   nil,
				repoGenreResult: true,
				reqDtoGenre:     dto.CreateGenreRequest{Name: "Genre 1", Image: &image1},
			},
			expected: expected{
				errType:      nil,
				valErrorMaps: nil,
			},
		},
		{
			name: "genreNotFound",
			fields: fields{
				repoGenreErr:    errs.NewNotFoundError("Genre", "Id", 1),
				repoUpdateErr:   nil,
				repoGenreResult: false,
				reqDtoGenre:     dto.CreateGenreRequest{Name: "Genre 1", Image: &image1},
			},
			expected: expected{
				errType:      errs.NewNotFoundError("Genre", "Id", 1),
				valErrorMaps: nil,
			},
		},
		{
			name: "repoUpdateError",
			fields: fields{
				repoGenreErr:    nil,
				repoUpdateErr:   errors.New("database failure"),
				repoGenreResult: true,
				reqDtoGenre:     dto.CreateGenreRequest{Name: "Genre 1", Image: &image1},
			},
			expected: expected{
				errType:      errors.New("database failure"),
				valErrorMaps: nil,
			},
		},
		{
			name: "validationErrors_RequiredName",
			fields: fields{
				reqDtoGenre: dto.CreateGenreRequest{Name: "", Image: &image1},
			},
			expected: expected{
				errType:      errValidation,
				valErrorMaps: map[string]string{"name": "Field is required"},
			},
		},
		{
			name: "validationErrors_RequiredImage",
			fields: fields{
				reqDtoGenre: dto.CreateGenreRequest{Name: "Genre 1", Image: nil},
			},
			expected: expected{
				errType:      errValidation,
				valErrorMaps: map[string]string{"image": "Field is required"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockGenreRepository)
			svc := NewGenreService(mockRepo, nil, nil, nil)
			ctx := context.Background()

			// Setup mock if applicable
			if tc.expected.valErrorMaps == nil {
				mockRepo.On("FindExistsGenreById", mock.Anything, 1).Return(tc.fields.repoGenreResult, tc.fields.repoGenreErr)

				if tc.fields.repoGenreErr == nil && tc.fields.repoGenreResult {
					// convert DTO to input model
					input := models.CreateGenreInput{
						Name:  tc.fields.reqDtoGenre.Name,
						Image: image1Bytes,
					}

					mockRepo.On("Update", ctx, input, 1).Return(tc.fields.repoUpdateErr)
				}
			}

			err := svc.UpdateGenre(ctx, tc.fields.reqDtoGenre, 1)

			if tc.expected.errType == nil {
				assert.NoError(t, err)
			} else if errors.Is(tc.expected.errType, errValidation) {
				var valErr *errs.BadRequestError
				assert.ErrorAs(t, err, &valErr)
				assert.Equal(t, tc.expected.valErrorMaps, valErr.Errors)
			} else {
				assert.Equal(t, tc.expected.errType, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenreService_DeleteGenre(t *testing.T) {
	type fields struct {
		repoDeleteErr   error
		repoGenreErr    error
		repoGenreResult bool
		genreId         int
	}

	type expected struct {
		err error
	}

	type scenario struct {
		name     string
		fields   fields
		expected expected
	}

	testCases := []scenario{
		{
			name: "success",
			fields: fields{
				repoDeleteErr:   nil,
				repoGenreResult: true,
				genreId:         1,
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "genreNotFound",
			fields: fields{
				repoDeleteErr:   nil,
				repoGenreErr:    errs.NewNotFoundError("Genre", "id", 1),
				repoGenreResult: false,
				genreId:         1,
			},
			expected: expected{
				err: errs.NewNotFoundError("Genre", "id", 1),
			},
		},
		{
			name: "deleteError",
			fields: fields{
				repoDeleteErr:   errors.New("database failure"),
				repoGenreErr:    nil,
				repoGenreResult: true,
				genreId:         1,
			},
			expected: expected{
				err: errors.New("database failure"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockGenreRepository)
			svc := NewGenreService(mockRepo, nil, nil, nil)

			ctx := context.Background()

			mockRepo.On("FindExistsGenreById", ctx, tc.fields.genreId).Return(tc.fields.repoGenreResult, tc.fields.repoGenreErr)

			if tc.fields.repoGenreErr == nil && tc.fields.repoGenreResult {
				mockRepo.On("Delete", ctx, tc.fields.genreId).Return(tc.fields.repoDeleteErr)
			}

			// Actual
			err := svc.DeleteGenre(ctx, tc.fields.genreId)

			// Assert
			if tc.expected.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expected.err, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
