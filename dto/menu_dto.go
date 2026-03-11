package dto

type MenuListRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type MenuResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type MenuListResponse struct {
	Total int64          `json:"total"`
	Data  []MenuResponse `json:"data"`
}

type CreateMenuRequest struct {
	Name string `json:"name" binding:"required"`
	Path string `json:"path" binding:"required"`
}

type UpdateMenuRequest struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

