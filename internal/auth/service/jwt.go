package service

import (
	"fmt"
	"time"

	e "github.com/BwezB/Wikno-backend/pkg/errors"
	

	"github.com/golang-jwt/jwt/v5"
)


// Claims struct for JWT
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func generateJWT(userID, email, secret string, expiery time.Duration) (string, error) {
	claims := Claims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiery)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", e.New("failed to generate token", ErrInternal, err)
    }

    return signedToken, nil
}

func verifyJWT(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Validate signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            return nil, e.New("invalid signing method", ErrInvalidToken, err)
        }
        return []byte(secret), nil
    })
	if err != nil {
        return nil, e.New("invalid token", ErrInvalidToken, err)
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, e.New("invalid token claims", ErrInvalidToken, nil)
}
