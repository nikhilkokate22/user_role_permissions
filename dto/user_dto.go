package dto

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	RoleID   uint   `json:"role_id"`
}

type UserListRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type UserListResponse struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Mobile    string `json:"mobile"`
	RoleID    uint   `json:"role_id"`
	CreatedAt string `json:"created_at"`
}

type UserListPaginatedResponse struct {
	TotalCount    int64              `json:"total_count"`
	FilteredCount int64              `json:"filtered_count"`
	Data          []UserListResponse `json:"data"`
}

type UpdateUserRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
}

type DeleteUserRequest struct {
	ID uint `json:"id" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	// Roles         []RoleResponse `json:"roles"`
	Token string `json:"token"`
}
