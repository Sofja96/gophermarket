package middleware

import (
	"fmt"
	"github.com/Sofja96/gophermarket.git/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	User string
}

const JwtSecret = "JWT_SECRET"

const TokenExp = time.Hour * 24

//const TokenExp = time.Second * 24

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
			//if token == "" {
			//	//return c.JSON(http.StatusUnauthorized, "You must be logged in to access this resource")
			//}
			log.Println(token, "header")

			//user, err := VerifyToken("Bearer" + token)
			log.Println(token, "token")
			user, err := VerifyToken(token)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusUnauthorized, "You must be logged in to access this resource")
			}

			//c.Request().Context().Value(models.ContextKeyUser)
			//log.Println(c.Request().Context().Value(models.ContextKeyUser))
			c.Set(models.ContextKeyUser, user)
			log.Println(c.Get(models.ContextKeyUser))
			//log.Println(c.Get(user))
			//log.Println(c.Get(token))
			//log.Println(c.Request().FormValue(models.ContextKeyUser))
			//	c.Get(models.ContextKeyUser)

			//claims.(jwt.MapClaims)["id"]
			//user := token.Claims.(*Token)
			//var newUser models.User
			//c.Request().SetBasicAuth(newUser.Login, newUser.Password)
			//username, _, _ := c.Request().BasicAuth()
			//log.Println(username, "basicAuth")

			if err = next(c); err != nil {
				c.Error(err)
			}

			return err

		}
	}
}
