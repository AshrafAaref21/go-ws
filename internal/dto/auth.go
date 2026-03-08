package dto

type EmailRegistrationRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type EmailLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshSessionRequest struct {
	RefreshToken string `json:"refresh_token"`
}
