package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/XRS0/ToTalkB/auth/db"
	"github.com/XRS0/ToTalkB/auth/middleware"
	"github.com/XRS0/ToTalkB/chat"
	"github.com/XRS0/ToTalkB/chat/pkg"
	"github.com/gin-gonic/gin"
)

func serveHome(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "home.html")
}

func main() {
	hub := chat.NewHub()
	go hub.Run()

	db, err := db.NewPostgresDB(db.Config{Host: "localhost", Port: "5432", Username: "postgres", Password: "postgres", DBName: "postgres", SSLMode: "disable"})
	if err != nil {
		log.Fatalf("failed to connect to db: %s\n", err.Error())
	}
	defer db.Close()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/chat/:chatId", serveHome)
	r.GET("/ws/:chatId", func(c *gin.Context) {
		userId := "1" // !!!
		chatId := c.Param("chatId")
		chat.ServeWs(hub, c.Writer, c.Request, db, userId, chatId)
	})
	api := r.Group("/api", middleware.UserIdentity)
	api.POST("/chat", func(c *gin.Context) {
		var input pkg.Chat

		uid, err := strconv.Atoi(c.GetString("userId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		input.OwnerId = uid

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		defer tx.Rollback()

		var chatID int
		err = db.QueryRowx(
			`INSERT INTO chats (name, created_by) 
         VALUES ($1, $2) 
         RETURNING id`,
			input.Name,
			input.OwnerId,
		).Scan(&chatID)
		if err != nil {
			log.Printf("Failed to create chat: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
			return
		}

		_, err = db.Exec(
			`INSERT INTO chat_members (user_id, chat_id) 
         VALUES ($1, $2)`,
			input.OwnerId,
			chatID,
		)
		if err != nil {
			log.Printf("Failed to add owner to chat members: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add chat member"})
			return
		}

		if err := tx.Commit(); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       chatID,
			"name":     input.Name,
			"owner_id": input.OwnerId,
		})
	})

	r.Run(":8081")
}
