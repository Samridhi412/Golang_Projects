package models

import "time"
type User struct{
	Username, PasswordHash, Role string
}
type TokenClaims struct{
	jwt.StandardClaims
	Role string `json:"role"`
	Csrf string `json:"csrf"`
}
const RefreshTokenValidTime = time.Hour * 72
const AuthTokenValidTime = time.Minute * 15

//once authtoken expires , frontend requests backend for new token using refreshtoken
func GenerateCSRFSecret()(string, error){
	return  randomstrings.GenerateRandomString(32)
}
