package dto

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type PaginationResult struct {
	TotalRecords int64       `json:"totalRecords"`
	Records      interface{} `json:"records"`
	CurrentPage  int         `json:"currentPage"`
	TotalPages   int         `json:"totalPages"`
}
