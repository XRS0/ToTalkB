package main

import (
	"fmt"
	"log"

	"github.com/XRS0/ToTalkB/auth"
	"github.com/XRS0/ToTalkB/auth/db"
	"github.com/XRS0/ToTalkB/auth/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.NewPostgresDB(db.Config{Host: "localhost", Port: "5432", Username: "postgres", Password: "postgres", DBName: "postgres", SSLMode: "disable"})
	if err != nil {
		log.Fatalf("Failed to connect to db: %s\n", err.Error())
	}
	defer db.Close()

	auth := &auth.Auth{DB: db}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.POST("/api/auth/sign-up", auth.SignUp)
	r.POST("/api/auth/sign-in", auth.SignIn)
	r.GET("/api/get-user", middleware.UserIdentity, auth.GetUser)

	fmt.Println("Auth Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %s\n", err.Error())
	}
}
