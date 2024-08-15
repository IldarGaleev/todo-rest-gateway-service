// Package httpdto contains HTTP DTO models
package httpdto

type LoginResponse struct {
	GeneralResponse
	Token string `json:"token"`
	// RefreshToken string `json:"refresh_token"`
}
