package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"

	"auth/db"
	"auth/pkg"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	salt = "x1n98r98y1xr2n8y"
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

	// валидация???

	id, err := uuid.NewRandom()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't to generate UUID: %s", err.Error())})
		return
	}

	hashedPassword := generatePasswordHash(input.Password)

	input.Id = id.String()
	input.Password = hashedPassword

	// такой чел уже есть???

	query := "INSERT INTO users (login, password_hash, name, role) values ($1, $2, $3, $4)"

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
