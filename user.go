package models

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

type UsersResponse struct {
	Users map[string]string `json:"users"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
