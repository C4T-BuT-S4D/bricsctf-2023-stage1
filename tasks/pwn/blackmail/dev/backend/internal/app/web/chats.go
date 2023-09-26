package web

import (
	"fmt"
	"net/http"

	"cbs.dev/brics/droidchat/internal/app"
	"cbs.dev/brics/droidchat/internal/app/config"

	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
	// Dependency injection
	Service app.ChatsService
}

func MountChats(apiRoot *gin.RouterGroup) *ChatRouter {
	r := ChatRouter{}
	apiRoot.Group("/chats").
		Use(RequireAuth).
		GET("/", r.GET_userChats).
		GET("/:id", r.GET_chat).
		POST("/:id", r.POST_message)
	return &r
}

func (h *ChatRouter) GET_userChats(c *gin.Context) {
	uid := GetUid(c)
	if chats, err := h.Service.GetChatsPreview(uid); err == nil {
		c.JSON(http.StatusOK, gin.H{"chats": chats})
	} else {
		panic(err)
	}
}

func (h *ChatRouter) GET_chat(c *gin.Context) {
	uidMe, uidOther, ok := h.getUids(c)
	if !ok {
		return
	}
	if chat, err := h.Service.GetChat(uidMe, uidOther); err == nil {
		c.JSON(http.StatusOK, gin.H{"chat": chat})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

type inputMessage struct {
	Type      string `json:"type" binding:"required,oneof=text sticker"`
	Text      string `json:"text" binding:"excluded_unless=Type text,required_if=Type text,lt=1024"`
	StickerId string `json:"sticker_id" binding:"excluded_unless=Type sticker,required_if=Type sticker"`
}

func (h *ChatRouter) POST_message(c *gin.Context) {
	uidMe, uidOther, ok := h.getUids(c)
	if !ok {
		return
	}

	var inMsg inputMessage
	if err := c.BindJSON(&inMsg); err != nil {
		fmt.Println(err)
		// !!!
		return
	}

	newMsg := app.Message{
		Type: inMsg.Type,
	}
	if inMsg.Type == app.MessageTypeSticker {
		stik, ok := config.AvailableStickers[inMsg.StickerId]
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "unknown sticker"})
			return
		}
		newMsg.Sticker = stik
	} else {
		newMsg.Text = &inMsg.Text
	}
	if err := h.Service.SendMessage(uidMe, uidOther, newMsg); err == nil {
		c.Status(http.StatusCreated)
	} else {
		c.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
	}
}

func (h *ChatRouter) getUids(c *gin.Context) (me app.Uid, other app.Uid, ok bool) {
	ok = false
	var idParam struct {
		Value int `uri:"id" binding:"required,numeric,gt=0"`
	}
	if err := c.BindUri(&idParam); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid id param"})
		return
	}

	me, other = GetUid(c), app.Uid(idParam.Value)

	if me == other {
		c.JSON(http.StatusForbidden, gin.H{"error": "amidst the whispers of the digital realm, an enigmatic hush befalls as your query seeks discourse with echoes only the mind can fathom."})
		return
	}

	ok = true
	return
}
