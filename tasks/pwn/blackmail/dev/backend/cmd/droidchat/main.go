package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"cbs.dev/brics/droidchat/internal/app"
	"cbs.dev/brics/droidchat/internal/app/playpen"
	"cbs.dev/brics/droidchat/internal/app/postgres"
	"cbs.dev/brics/droidchat/internal/app/web"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	time.Sleep(5*time.Second)
	db, err := sql.Open("postgres", os.Getenv("DB_URI"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(32)
//	db.SetMaxIdleConns(16)

	usersSvc, chatsSvc := setupServices(db)

	engine := gin.Default()
	appRouter := web.NewRouter(os.Getenv("STATIC_HOST"), engine.Group("/api"))
	engine.Static("/static", "./static")

	appRouter.Users.Service = usersSvc
	appRouter.Chats.Service = chatsSvc

	if err := engine.Run("0.0.0.0:3000"); err != nil {
		log.Fatal(err)
	}
}

func setupServices(db *sql.DB) (app.UserService, app.ChatsService) {
	users := postgres.NewUserService(db)
	chats := postgres.NewChatsService(db, users)
	if _, ok := os.LookupEnv("NOPLAYPEN"); ok {
		log.Println("NOT ISOLATING PLAYERS!!!")
		return users, chats
	}
	return users, playpen.NewChatsService(users, chats)
}
