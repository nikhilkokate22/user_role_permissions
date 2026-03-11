package dto

type CreateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type ListRoleRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UpdateRoleRequest struct {
	ID   uint   `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DeleteRoleRequest struct {
	ID uint `json:"id" binding:"required"`
}

type AssignMenuPermissionRequest struct {
	Id    int          `json:"id"`
	Menus []AssignMenu `json:"menus"`
}
type AssignMenuPermissionResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type AssignMenu struct {
	Id          int   `json:"menu_id"`
	Permissions []int `json:"permissions"`
}