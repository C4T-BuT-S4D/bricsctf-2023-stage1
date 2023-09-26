package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	apps "cbs.dev/brics/droidchat/runner/internal/db"
	"github.com/gin-gonic/gin"
)

type web struct {
	db *apps.AppsDB
}

func InitWeb(db *apps.AppsDB) (*gin.Engine, *web) {
	w := &web{db: db}

	e := gin.Default()

	e.LoadHTMLGlob("templates/*")
	e.Static("/static", "./static")

	e.MaxMultipartMemory = 6 << 20

	e.SetTrustedProxies(nil)
	e.TrustedPlatform = "CF-Connecting-IP"

	e.GET("/", w.GET_mainPage)
	e.POST("/", w.POST_submitApp)

	return e, w
}

func (w *web) GET_mainPage(c *gin.Context) {
	ip := c.ClientIP()

	sub, err := w.db.GetLastIPSubmission(ip)
	if err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		log.Printf("get last sub: %v\n", err)
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", sub)
}

func (w *web) POST_submitApp(c *gin.Context) {
	ip := c.ClientIP()

	t := time.Now()
	if lt, err := w.db.GetLastAppTime(ip); err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		log.Printf("get last api time err: %v\n", err)
		return
	} else if t.Sub(lt) <= 5*time.Minute {
		c.String(http.StatusTooManyRequests,
			fmt.Sprintf("you must wait for another %v", (5*time.Minute - t.Sub(lt).Abs()).String()))
		return
	}

	var apk []byte
	{
		hdr, _ := c.FormFile("apk")
		if hdr == nil {
			c.String(http.StatusBadRequest, "apk upload missing")
			return
		}
		if hdr.Size > 4<<20 {
			c.String(http.StatusUnprocessableEntity, "apk must be under 4 MiB")
			return
		}
		f, _ := hdr.Open()
		apkBytes, _ := io.ReadAll(f)
		apk = apkBytes
	}

	if _, err := w.db.Submit(ip, apk); err == apps.ErrDuplicateApk {
		c.String(http.StatusConflict, "this apk was already submitted")
	} else if err != nil {
		c.String(http.StatusInternalServerError, "internal error")
		log.Printf("get last api time err: %v\n", err)
	} else {
		// JS script will highlight the new submission
		c.Redirect(http.StatusFound, fmt.Sprintf("/"))
	}
}

