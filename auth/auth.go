package auth

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/XRS0/ToTalkB/auth/pkg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	salt       = "x1n98r98y1xr2n8y"
	signingKey = "cenuc2enuc2ni92c"
	accessTTL  = 7 * 24 * time.Hour
)

type Auth struct {
	DB *sqlx.DB
}

func (a *Auth) SignUp(c *gin.Context) {
	var input pkg.User

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	input.Password = generatePasswordHash(input.Password)

	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(login) = LOWER($1))"
	err := a.DB.QueryRow(query, input.Login).Scan(&exists)
	if err != nil || exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "login is already taken"})
		return
	}

	query = "INSERT INTO users (login, password_hash, name, role) values ($1, $2, $3, $4)"

	_, err = a.DB.Exec(query, input.Login, input.Password, input.Name, input.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create user in db: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *Auth) SignIn(c *gin.Context) {
	var input pkg.SignInRequest

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	hashedPassword := generatePasswordHash(input.Password)
	input.Password = hashedPassword

	var userId string

	query := "SELECT id FROM users WHERE login = $1 AND password_hash = $2"
	err := a.DB.Get(&userId, query, input.Login, input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": fmt.Sprintf("can't get user id: %s", err.Error())})
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(accessTTL).Unix(),
		Subject:   userId,
	})

	tokenString, err := accessToken.SignedString([]byte(signingKey))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": fmt.Sprintf("can't create new access token: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func ParseAccessToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
					return claims.Subject, errors.New("token has expired")
				}
			}
		}
		return "", err
	}

	if !token.Valid {
		return "", errors.New("token is invalid")
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", errors.New("token claims are not of type *TokenClaims")
	}

	return claims.Subject, nil
}
