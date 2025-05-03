package dto

type Image struct {
	Src      string `json:"src"`
	BlurHash string `json:"blur_hash"`
}

type Pagination struct {
	Total    int `json:"total"`
	PageSize int `json:"page_size"`
	Page     int `json:"page"`
}
