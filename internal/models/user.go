package models

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`

	Dob string `json:"dob" validate:"required,datetime=02-01-2006"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=100"`

	Dob string `json:"dob" validate:"required,datetime=02-01-2006"`
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
}

type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Dob  string `json:"dob"`
	Age  int    `json:"age"`
}

type PaginatedUsersResponse struct {
	Data       []UserWithAgeResponse `json:"data"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalItems int64                 `json:"total_items"`
	TotalPages int                   `json:"total_pages"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
