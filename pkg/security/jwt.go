package security

import "github.com/golang-jwt/jwt/v5"

type Jwt struct {
	SecretKey     []byte
	Claims        jwt.Claims
	SigningMethod jwt.SigningMethod
}

func (j *Jwt) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(j.SigningMethod, j.Claims)
	return token.SignedString(j.SecretKey)
}
