package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	_ "github.com/lib/pq"

	apps "cbs.dev/brics/droidchat/runner/internal/db"
)

func main() {
	time.Sleep(3*time.Second)
	db, err := sql.Open("postgres", requireEnv("DB_URI"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	appsDb := apps.NewAppsDB(db)

	emu := NewEmulator(os.Args[1])

	log.Println("starting main loop")
	go func() {
		for {
			time.Sleep(5*time.Second)
			ps, err := appsDb.GetPendingSubmission()
			if err != nil {
				log.Fatal(err)  // or swallow it?
			}
			if ps == nil {
				// log.Println("no new submissions")
				continue
			}
			
			log.Printf("new app %s\n", ps.Id)
			path, err := dlApk(ps)
			if err != nil {
				// ???
				log.Printf("error saving apk: %v\n", err)
				appsDb.SetStatus(ps.Id, false, "sorry, try again")
				continue
			}

			if err := checkApp(path); err != nil {
				log.Println(err)				
				appsDb.SetStatus(ps.Id, false, err.Error())
				continue
			}

			if err := emu.InstallApp(path); err != nil {
				log.Println(err)				
				appsDb.SetStatus(ps.Id, false, "invalid apk")
				continue
			}
			// if err := emu.StartApp("ru.bricsctf.droidchat"); err != nil {
			// 	log.Println(err)
			// 	appsDb.SetStatus(ps.Id, false, "error starting app")
			// 	if err := emu.Reset(); err != nil {
			// 		log.Fatal(err)
			// 	}
			// 	continue
			// }
			if err := emu.StartApp("ru.superappstore.newapp"); err != nil {
				log.Println(err)
				appsDb.SetStatus(ps.Id, false, "error starting app")
				if err := emu.Reset(); err != nil {
					log.Fatal(err)
				}
				continue
			}

			os.Remove(path)

			time.Sleep(5*time.Second)

			log.Println("time's out, resetting emulator")
			if err := emu.Reset(); err != nil {
				log.Fatal(err)
			}

			appsDb.SetStatus(ps.Id, false, "review passed")
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func dlApk(ps *apps.PendingSubmission) (string, error) {
	path := fmt.Sprintf("%s/%s.apk", requireEnv("APPS_DIR"), ps.Id)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	_, err = f.Write(ps.Apk)
	return path, err
}

func requireEnv(key string) string {
	if v, ok := os.LookupEnv(key); !ok {
		log.Fatalf("%v env missing", v)
		return ""
	} else {
		return v
	}
}

func checkApp(path string) error {
	cmd := exec.Command("apkanalyzer", "manifest", "permissions", path)

	if b, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return fmt.Errorf("invalid apk")
	} else if strings.Contains(string(b), "android.permission.INTERNET") {
		return fmt.Errorf("app requires additional checks")
	}

	cmd = exec.Command("apkanalyzer", "files", "list", path)
	if b, err := cmd.Output(); err != nil {
		log.Println(string(err.(*exec.ExitError).Stderr))
		return fmt.Errorf("invalid apk")
	} else if strings.Contains(string(b), "/lib/") {
		return fmt.Errorf("app requires additional checks")
	}
	
	return nil
}
