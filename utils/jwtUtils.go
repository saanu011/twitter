package utils

import "github.com/form3tech-oss/jwt-go"

// UserClaims contains information currently being sent in a JWT
type UserClaims struct {
	Application       string `json:"application"`
	ID                string `json:"_id"`
	Username          string `json:"username"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	LastLogin         string `json:"last_login"`
	LastPasswordReset string `json:"last_password_reset"`
	jwt.StandardClaims
}
