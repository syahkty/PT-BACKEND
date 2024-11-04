package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret")

type JWTAuth struct{}

func NewJWTAuth() *JWTAuth {
    return &JWTAuth{}
}

type Claims struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username"`
    Role     string `json:"role"`

    jwt.StandardClaims
}

func (a *JWTAuth) GenerateJWT(username, role string,id uint) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        ID : id,
        Username: username,
        Role:     role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func (a *JWTAuth) ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }
    fmt.Printf("Token validated. Role: %v\n", claims.Role)
    return claims, nil
}
