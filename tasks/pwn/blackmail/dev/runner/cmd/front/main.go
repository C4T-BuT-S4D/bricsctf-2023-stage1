package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	apps "cbs.dev/brics/droidchat/runner/internal/db"
	_ "github.com/lib/pq"
)

const ENV_DB_URI = "DB_URI"

func main() {
	time.Sleep(3*time.Second)
	db := getDB()
	defer db.Close()
	appsDb := apps.NewAppsDB(db)

	if err := appsDb.EnsureSchema(); err != nil {
		log.Fatal(err)
	}

	router, _ := InitWeb(appsDb)
	if err := router.Run("0.0.0.0:3000"); err != nil {
		log.Fatal(err)
	}
}

func requireEnv(key string) string {
	if v, ok := os.LookupEnv(key); !ok {
		log.Fatalf("%v env missing", v)
		return ""
	} else {
		return v
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("postgres", requireEnv(ENV_DB_URI))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}
