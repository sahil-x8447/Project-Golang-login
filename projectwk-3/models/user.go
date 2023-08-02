package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GetUserDetailsRequest represents the request structure for GetUserDetailsHandler.
type GetUserDetailsRequest struct {
	Username string `json:"username"`
}

// GetUserDetailsResponse represents the response structure for GetUserDetailsHandler.
type GetUserDetailsResponse struct {
	Username string `json:"username"`
}

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
	Token   string `json:"token"`
}

type UsersResponse struct {
	Users map[string]string `json:"users"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
