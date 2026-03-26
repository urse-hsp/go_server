package dto

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

type PageRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type PageResponse struct {
	List     any `json:"list"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
