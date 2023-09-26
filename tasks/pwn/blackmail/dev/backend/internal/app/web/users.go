package web

import (
	"fmt"
	"net/http"

	"cbs.dev/brics/droidchat/internal/app"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	// Dependency injection
	Service app.UserService
}

func MountUsers(apiRoot *gin.RouterGroup) *UserRouter {
	r := UserRouter{}
	apiRoot.Group("/users").
		POST("/", r.POST_createUser).
		GET("/", RequireAuth, r.GET_allUsers).
		POST("/token", r.POST_getToken).
		GET("/:id", RequireAuth, r.GET_user)
	return &r
}

func (h *UserRouter) POST_createUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=4,max=72,alphanum"`
		Password string `json:"password" binding:"required,min=8,max=72"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	fmt.Println(h.Service)

	if newUser, err := h.Service.Register(req.Username, req.Password); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusCreated, gin.H{"user": newUser})
	}
}

func (h *UserRouter) POST_getToken(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=4,max=72"`
		Password string `json:"password" binding:"required,min=8,max=72"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	if user, err := h.Service.GetByCreds(req.Username, req.Password); err == nil {
		c.JSON(http.StatusOK, gin.H{"token": MakeToken(user.Id)})
	} else {
		// or 404?
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong username/password"})
	}
}

func (h *UserRouter) GET_user(c *gin.Context) {
	var id app.Uid
	{
		var idParam struct {
			Me string `uri:"id" binding:"eq=me"`
			Id int    `uri:"id" binding:"numeric"`
		}
		c.ShouldBindUri(&idParam)

		if idParam.Me == "me" {
			id = GetUid(c)
		} else if idParam.Id > 0 {
			id = app.Uid(idParam.Id)
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "id path param invalid"})
			return
		}
	}

	if user, err := h.Service.GetById(id); err == nil {
		c.JSON(http.StatusOK, gin.H{"user": user})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}
}

func (h *UserRouter) GET_allUsers(c *gin.Context) {
	if users, err := h.Service.GetBots(); err == nil {
		c.JSON(http.StatusOK, gin.H{"users": users})
	} else {
		// ?
		c.JSON(http.StatusTeapot, gin.H{"error": err.Error()})
	}
}
