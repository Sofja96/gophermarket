package middleware

import (
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	User string
}

const JwtSecret = "JWT_SECRET"

const TokenExp = time.Hour * 24

func CreateToken(user string) (string, error) {
	claims := Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(JwtSecret), nil
		})
	if err != nil {
		return "", fmt.Errorf("error on parsing token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	fmt.Println("Token is valid")
	return claims.User, nil
}

const BearerSchema = "Bearer "

func ValidateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, "You must be logged in to access this resource")
			}
			token := authHeader[len(BearerSchema):]

			user, err := VerifyToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, "You must be logged in to access this resource")
			}

			c.Set(models.ContextKeyUser, user)

			if err = next(c); err != nil {
				c.Error(err)
			}

			return err

		}
	}
}
