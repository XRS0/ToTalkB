package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	salt = "x1n98r98y1xr2n8y"
)

var users = map[string]*User{}

type User struct {
	Id       string `json:"-" db:"id"`
	Login    string `json:"login" db:"login" binding:"required"`
	Password string `json:"password" db:"password_hash" binding:"required"`
	Name     string `json:"name" db:"name" binding:"required"`
	Role     string `json:"role" db:"role" binding:"required"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func main() {
	r := gin.Default()

	r.POST("/sign-up", signUp)
	r.GET("/users", getAllUsers)
	r.GET("/users/:id", getUserByID)

	fmt.Println("Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func getAllUsers(c *gin.Context) {
	userList := make(map[string]UserResponse)
	for id, user := range users {
		userList[id] = UserResponse{
			Name:  user.Name,
			Role:  user.Role,
		}
	}
	c.JSON(http.StatusOK, gin.H{"users": userList})
}

func getUserByID(c *gin.Context) {
	id := c.Param("id")
	user, exists := users[id]
	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	response := UserResponse{
		Name:  user.Name,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, response)
}

func signUp(c *gin.Context) {
	var input User

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": fmt.Sprintf("can't bind JSON: %s", err.Error())})
		return
	}

	// валидация???

	if err := createUser(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": fmt.Sprintf("can't create user: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func createUser(input *User) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %s", err.Error())
	}

	hashedPassword := generatePasswordHash(input.Password)

	input.Id = id.String()
	input.Password = hashedPassword

	users[input.Id] = input

	return nil
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
