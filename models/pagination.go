package models

type Pagination struct {
	PerPage    int `json:"perPage"`
	Page       int `json:"page"`
	TotalPages int `json:"totalPages"`
	TotalCount int `json:"totalCount"`
}
