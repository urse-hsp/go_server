package v1

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
}

type PageRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"pageSize" binding:"required,min=1,max=100"`
}

type PageSizeResponse struct {
	Total    int `json:"total"`    // 总数
	Page     int `json:"page"`     // 页码
	PageSize int `json:"pageSize"` // 条数
}

type PageResponse[T any] struct {
	Data []T `json:"data"` // 列表
	PageSizeResponse
}
