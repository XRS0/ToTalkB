package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"time"

	"auth/db"
	"auth/pkg"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	salt       = "x1n98r98y1xr2n8y"
	signingKey = "cenuc2enuc2ni92c"
	accessTTL  = 7 * 24 * time.Hour
)

type Auth struct {
	db *sqlx.DB
}

func main() {
	db, err := db.NewPostgresDB(db.Config{Host: "localhost", Port: "5432", Username: "postgres", Password: "postgres", DBName: "postgres", SSLMode: "disable"})
	if err != nil {
		log.Fatalf("Failed to connect to db: %s\n", err.Error())
	}
	defer db.Close()

	auth := &Auth{db: db}

	r := gin.Default()
	r.POST("/api/auth/sign-up", auth.signUp)
	r.POST("/api/auth/sign-in", auth.signIn)

	fmt.Println("Auth Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}

func (a *Auth) signUp(c *gin.Context) {
	var input pkg.User

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't to generate UUID: %s", err.Error())})
		return
	}

	hashedPassword := generatePasswordHash(input.Password)

	input.Id = id.String()
	input.Password = hashedPassword

	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(login) = LOWER($1))"
	err = a.db.QueryRow(query, input.Login).Scan(&exists)
	if err != nil || exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "login is already taken"})
		return
	}

	query = "INSERT INTO users (login, password_hash, name, role) values ($1, $2, $3, $4)"

	_, err = a.db.Exec(query, input.Login, input.Password, input.Name, input.Role)
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

func (a *Auth) signIn(c *gin.Context) {
	var input pkg.SignInRequest

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	hashedPassword := generatePasswordHash(input.Password)
	input.Password = hashedPassword

	var userId string

	query := "SELECT id FROM users WHERE login = $1 AND password_hash = $2"
	err := a.db.Get(&userId, query, input.Login, input.Password)
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
