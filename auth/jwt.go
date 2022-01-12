package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// secret key being used to sign tokens
var (
	SecretKey = []byte("lessbutterhabitapi")
)

func GenerateToken(email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", errors.New("generate jwt error")
	}
	return tokenString, nil
}

// Middleware에서 parsing해서 valid한지 체크 + Refreshtoken요청할때 ParseToken으로 사용됨
func ParseToken(tokenStr string) (string, error) {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["userid"].(string)
		return userID, nil
	} else {
		return "", errors.New("token expired or incorrect")
	}
}
