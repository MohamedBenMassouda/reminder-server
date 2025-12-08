package models

type Pagination struct {
	Offset int `json:"offset" form:"offset"`
	Limit  int `json:"limit" form:"limit"`
	Total  int `json:"total"`
	Items  any `json:"items"`
}
