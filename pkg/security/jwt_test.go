package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"testing"
)

func TestJwt_GenerateToken(t *testing.T) {
	key := os.Getenv("JWT_SECRET_KEY")
	JwtStruct := Jwt{
		SecretKey: []byte(key),
		Claims: jwt.MapClaims{
			"username": "yofs",
		},
		SigningMethod: jwt.SigningMethodHS256,
	}

	token, err := JwtStruct.GenerateToken()
	if err != nil {
		t.Fatalf("An error has occurred while trying to generate the jwt token %v", err)
	}
	fmt.Println(token)
}
