package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"time"

	"cbs.dev/brics/droidchat/internal/app"
	"cbs.dev/brics/droidchat/internal/app/postgres"
)

var (
	users app.UserService
	chats app.ChatsService
	db    *sql.DB
)

func main() {
	time.Sleep(3*time.Second)  // Wait for db init
	var err error
	db, err = sql.Open("postgres", os.Getenv("DB_URI"))
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	//db.SetMaxIdleConns(14)
	db.SetMaxOpenConns(16)
	// db.SetConnMaxIdleTime(6*time.Second)
	// db.SetConnMaxLifetime(10*time.Second)

	users = postgres.NewUserService(db)
	chats = postgres.NewChatsService(db, users)

	var admin *app.User
	_, err = users.GetByName("admin")
	if err == app.ErrNotFound {
		admin, err = users.Register("admin", thePassword+"iamadm1n")
		if err != nil {
			log.Fatalf("can't register admin: %v", err)
		}
	} else if err != nil && err != app.ErrNotFound {
		log.Fatalf("error checking for admin: %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, u := range usernames {
		var user *app.User
		user, err := users.GetByName(u)
		if err == app.ErrNotFound {
			user, err = users.Register(u, thePassword)
			if err != nil {
				log.Printf("can't register bot %v\n", err)
				return
			}
			if err := makeBot(db, user.Id); err != nil {
				log.Printf("can't make a bot %v\n", err)
				return
			}
		} else if err != nil {
			log.Printf("can't use username: %v\n", err)
			return
		}
		go bot(ctx, user)
	}

	if admin != nil {
		if err := sendAdminMsgs(admin.Id); err != nil {
			log.Fatalf("error planting flag: %v\n", err)
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

// randomize
func sendAdminMsgs(adminUid app.Uid) error {
	flagMsg := text("Keep this secret please! " + os.Getenv("FLAG"))

	users, err := users.GetBots()
	if err != nil {
		return err
	}

	return chats.SendMessage(adminUid, users[2].Id, flagMsg)
}

func bot(ctx context.Context, u *app.User) {
	log := log.New(log.Default().Writer(), "<"+u.Username+"> ", log.LstdFlags|log.Lmsgprefix)

	log.Println("ready")
	for {
		time.Sleep(2 * time.Second)
		previews, err := chats.GetChatsPreview(u.Id)
		// log.Printf("user: %#v\n previews: %v", u, len(previews))
		if err != nil {
			log.Printf("error fetching new msgs: %v\n", err)
			continue
		}

		cnt := 0
		for _, c := range previews {
			if c.Messages[0].From != u.Id {
				chats.SendMessage(u.Id, c.With, randomMsg())
				cnt++
			}
		}
		if cnt > 0 {
			log.Printf("responded to %d msgs\n", cnt)
		}

		select {
		case <-ctx.Done():
			log.Println("stopping")
			return
		default:
		}
	}
}

func makeBot(db *sql.DB, id app.Uid) error {
	_, err := db.Query(`UPDATE "user" SET is_bot=TRUE WHERE id=$1`, id)
	return err
}
