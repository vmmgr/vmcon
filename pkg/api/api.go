package api

import (
	"github.com/gin-gonic/gin"
	controller "github.com/vmmgr/controller/pkg/api/core/controller/v0"
	group "github.com/vmmgr/controller/pkg/api/core/group/v0"
	notice "github.com/vmmgr/controller/pkg/api/core/notice/v0"
	ticket "github.com/vmmgr/controller/pkg/api/core/support/ticket/v0"
	token "github.com/vmmgr/controller/pkg/api/core/token/v0"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	user "github.com/vmmgr/controller/pkg/api/core/user/v0"
	"log"
	"net/http"
	"strconv"
)

func AdminRestAPI() error {

	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatAdmin)

			// Notice
			//
			v1.POST("/notice", notice.AddAdmin)
			v1.DELETE("/notice/:id", notice.DeleteAdmin)
			v1.GET("/notice", notice.GetAllAdmin)
			v1.GET("/notice/:id", notice.GetAdmin)
			v1.PUT("/notice/:id", notice.UpdateAdmin)

			//
			// User
			//
			// User Create
			v1.POST("/user", user.AddAdmin)
			// User Delete
			v1.DELETE("/user", user.DeleteAdmin)
			// User Update
			v1.PUT("/user", user.UpdateAdmin)
			v1.GET("/user", user.GetAdmin)
			v1.GET("/user/:id", user.GetAdmin)
			//
			// Token
			//
			v1.POST("/token/generate", token.GenerateAdmin)

			v1.POST("/token", token.GenerateAdmin)
			// Token Delete
			v1.DELETE("/token", token.Delete)
			v1.DELETE("/token/:id", token.DeleteAdmin)
			// Token Update
			v1.PUT("/token", token.UpdateAdmin)
			v1.GET("/token", token.GetAllAdmin)
			v1.GET("/token/:id", token.GetAdmin)
			//
			// Group
			//
			v1.POST("/group", group.AddAdmin)
			// Group Delete
			v1.DELETE("/group", group.DeleteAdmin)
			// Group Update
			v1.PUT("/group", group.UpdateAdmin)
			v1.GET("/group", group.GetAllAdmin)
			v1.GET("/group/:id", group.GetAdmin)

			//
			// Support
			//
			v1.POST("/support", ticket.CreateAdmin)
			v1.GET("/support", ticket.GetAllAdmin)
			//v1.POST("/support/:id", chat.AddAdmin)
			v1.GET("/support/:id", ticket.GetAdmin)
			v1.PUT("/support/:id", ticket.UpdateAdmin)
		}
	}
	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetAdminWebSocket)
		}
	}

	go ticket.HandleMessagesAdmin()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Admin.Port), router))
	return nil
}

func UserRestAPI() {
	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatUser)

			//
			// User
			//
			// User Delete
			//router.DELETE("/user", user.Delete)
			// User Get
			v1.GET("/user", user.Get)
			v1.GET("/user/all", user.GetGroup)
			// User ID Get
			// v1.GET("/user/:id",user.GetId)
			// User Update
			v1.PUT("/user/:id", user.Update)
			// User Mail MailVerify
			v1.GET("/user/verify/:token", user.MailVerify)
			//
			// Token
			//
			// get token for CHAP authentication
			v1.GET("/token/init", token.GenerateInit)
			// get token for user
			v1.GET("/token", token.Generate)
			// delete
			v1.DELETE("/token", token.Delete)
			//
			// Group
			//
			// Group Create
			v1.POST("/group", group.Add)
			v1.GET("/group", group.Get)
			v1.PUT("/group", group.Update)
			v1.GET("/group/all", group.GetAll)

			//
			// Support
			//
			v1.POST("/support", ticket.Create)
			v1.GET("/support", ticket.GetTitle)
			v1.GET("/support/:id", ticket.Get)
			//
			// Notice
			//
			v1.GET("/notice", notice.Get)

			// 現在検討中

			// Network JPNIC Admin
			//v1.POST("/group/network/jpnic/admin", jpnicAdmin.Add)
			//v1.DELETE("/group/network/jpnic/admin", jpnicAdmin.Delete)
			//v1.GET("/group/network/jpnic/admin", jpnicAdmin.Get)
			// Network JPNIC Tech
			//v1.POST("/group/network/jpnic/tech", jpnicTech.Add)
			//v1.DELETE("/group/network/jpnic/tech", jpnicTech.Delete)
			//v1.GET("/group/network/jpnic/tech", jpnicTech.Get)
		}
	}

	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetWebSocket)
		}
	}

	go ticket.HandleMessages()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.User.Port), router))
}

func cors(c *gin.Context) {

	//c.Header("Access-Control-Allow-Headers", "Accept, Content-ID, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-ID", "application/json")
	c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}