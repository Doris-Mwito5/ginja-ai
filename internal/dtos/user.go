package dtos

// RegisterRequest is the inbound payload for POST /auth/register
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest is the inbound payload for POST /login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is the response for POST /login
type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	User        interface{} `json:"user"`
	ExpiresAt   string      `json:"expires_at"`
}
