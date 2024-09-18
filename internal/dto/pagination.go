package dto

type Pagination struct {
	Page  int `json:"current_page"`
	Count int `json:"page_count"`
	Total int `json:"total"`
}
