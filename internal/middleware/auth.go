package middleware

import (
	"fmt"
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

	fmt.Println("Token os valid")
	return claims.User, nil
}

//func (a *Authenticator) RegisterUser(ctx context.Context, user, password string) (string, error) {
//	if !checkCredentials(user, password) {
//		return "", authenticatorer.ErrEmptyCredentials
//	}
//
//	exists, err := a.storage.UserExists(ctx, user)
//	if err != nil {
//		return "", fmt.Errorf("error on checking user existance: %w", err)
//	}
//	if exists {
//		return "", authenticatorer.ErrUserExists
//	}
//
//	hash, salt, err := generateHashAndSalt(user, password)
//	if err != nil {
//		return "", fmt.Errorf("error on generating user hash and salt: %w", err)
//	}
//
//	if err := a.storage.RegisterUser(ctx, user, hash, salt); err != nil {
//		return "", fmt.Errorf("error on registering user in storage: %w", err)
//	}
//
//	token, err := a.generateToken(user)
//	if err != nil {
//		return "", fmt.Errorf("error on generating token: %w", err)
//	}
//
//	return "Bearer " + token, nil
//}

func ValidateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			//if len(user) == 0 && len(password) == 0 {
			//	return c.JSON(http.StatusBadRequest, "empty credentials")
			//}
			//token, err := CreateToken(user)
			//if err != nil {
			//	return c.JSON(http.StatusInternalServerError, "error on generating token")
			//}
			token := c.Request().Header.Get("Authorization")
			if err != nil {
				c.JSON(http.StatusUnauthorized, "You must be logged in to access this resource")
			}
			_, err = VerifyToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, "You must be logged in to access this resource")
			}

			if err = next(c); err != nil {
				c.Error(err)
			}

			return err

		}
	}
}
