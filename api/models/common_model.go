package models

type PaginationQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type PaginationResponse struct {
	Total int64         `json:"total"`
	List  []interface{} `json:"list"`
}
