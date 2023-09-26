package web

import (
	"fmt"
	"net/http"

	"cbs.dev/brics/droidchat/internal/app"
	"cbs.dev/brics/droidchat/internal/app/config"
	"github.com/gin-gonic/gin"
)

var staticUrl string
var stickers []*app.Sticker

func MountStickers(hostname string, apiRoot *gin.RouterGroup) {
	staticUrl = fmt.Sprintf("https://%s/static/", hostname)
	apiRoot.GET("/stickers/", RequireAuth, GET_allStickers)

	for _, sticker := range config.AvailableStickers {
		sticker.Url = getStickerUrl(sticker.Url)
		stickers = append(stickers, sticker)
	}
}

func GET_allStickers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"stickers": stickers})
}

func getStickerUrl(filename string) string {
	return staticUrl + filename
}
