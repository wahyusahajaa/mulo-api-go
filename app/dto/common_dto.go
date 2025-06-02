package dto

type Image struct {
	Src      string `json:"src"`
	BlurHash string `json:"blur_hash"`
} //@name Image

type Pagination struct {
	Total    int `json:"total"`
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
} //@name Pagination

type ResponseMessage struct {
	Message string `json:"message"`
} // @name ResponseMessage

// ResponseWithPagination[T,X]
// @Description Response data with pagination
type ResponseWithPagination[T any, X any] struct {
	Data       T `json:"data"`
	Pagination X `json:"pagination"`
} // @name ResponseWithPagination

// ResponseWithPagination[T]
// @Description Response with single data
type ResponseWithData[T any] struct {
	Data T `json:"data"`
} // @name ResponseWithData

// ResponseToken[T,X]
// @Description Response with token
type ResponseWithToken[T any, X any] struct {
	Message     T `json:"message"`
	AccessToken X `json:"access_token"`
} // @name ResponseWithData

type ResponseErrorWithDetails struct {
	ErrorResponse
	Errors map[string]string
} //@name ResponseErrorWithDetails

type ResponseErrorWithRequestId struct {
	ErrorResponse
	RequestId string
} //@name ResponseErrorWithRequestId

type ErrorResponse struct {
	Message string `json:"message"`
} //@name ErrorResponse

type ValidationErrorResponse struct {
	Message string            `json:"message" example:"Validation failed"`
	Errors  map[string]string `json:"errors,omitempty" swaggertype:"object"` // optional
} // @name ValidationErrorResponse

type InternalErrorResponse struct {
	Message   string `json:"message" example:"Internal server error"`
	RequestId string `json:"request_id" example:"abcd-1234"`
} //@name InternalErrorResponse
