package dto

type Pagination struct {
	Page  int `json:"current_page"`
	Count int `json:"page_count"`
	Total int `json:"total"`
}

func NewPagination(page, count, total int) *Pagination {
	return &Pagination{
		Page:  page,
		Count: count,
		Total: total,
	}
}
