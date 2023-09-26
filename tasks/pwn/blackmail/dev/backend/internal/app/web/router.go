package web

import "github.com/gin-gonic/gin"

type Router struct {
	Users *UserRouter
	Chats *ChatRouter
}

func NewRouter(hostname string, rootGroup *gin.RouterGroup) Router {
	r := Router{
		Users: MountUsers(rootGroup),
		Chats: MountChats(rootGroup),
	}
	MountStickers(hostname, rootGroup)
	return r
}
