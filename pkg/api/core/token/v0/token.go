package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/core/tool/hash"
	logging "github.com/vmmgr/controller/pkg/api/core/tool/log"
	toolToken "github.com/vmmgr/controller/pkg/api/core/tool/token"
	"github.com/vmmgr/controller/pkg/api/core/user"
	dbToken "github.com/vmmgr/controller/pkg/api/store/token/v0"
	dbUser "github.com/vmmgr/controller/pkg/api/store/user/v0"
	"log"
	"net/http"
	"strings"
	"time"
)

func GenerateInit(c *gin.Context) {
	ip := c.ClientIP()
	userToken := c.Request.Header.Get("USER_TOKEN")
	log.Println("userToken: " + userToken)
	tmpToken, _ := toolToken.Generate(2)
	err := dbToken.Create(&token.Token{ExpiredAt: time.Now().Add(30 * time.Minute), UID: 0, Status: 0,
		UserToken: userToken, TmpToken: tmpToken, Debug: ip, Admin: false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
	} else {
		log.Println("Time: " + time.Now().String() + " IP: " + c.ClientIP())
		logging.WriteLog("IP: " + c.ClientIP())
		c.JSON(http.StatusOK, &token.ResultTmpToken{Status: true, Token: tmpToken})
	}
}

func Generate(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	hashPass := c.Request.Header.Get("HASH_PASS")
	mail := c.Request.Header.Get("Email")
	tokenResult := dbToken.Get(token.UserToken, &token.Token{UserToken: userToken})
	if tokenResult.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: tokenResult.Err.Error()})
		return
	}
	userResult := dbUser.Get(user.Email, &user.User{Email: mail})
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: userResult.Err.Error()})
		return
	}

	if !userResult.User[0].MailVerify {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: fmt.Sprintf("You don't have email verification.")})
		return
	}

	if userResult.User[0].Status >= 100 {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: fmt.Sprintf("status error")})
		return
	}

	if hash.Generate(userResult.User[0].Pass+tokenResult.Token[0].TmpToken) != strings.ToUpper(hashPass) {
		log.Println(userResult.User[0].Pass)
		log.Println(tokenResult.Token[0].TmpToken)
		log.Println("hash(server): " + hash.Generate(userResult.User[0].Pass+tokenResult.Token[0].TmpToken))
		log.Println("hash(client): " + hashPass)
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: "not match"})
		return
	}
	accessToken, _ := toolToken.Generate(2)
	err := dbToken.Update(token.AddToken, &token.Token{Model: gorm.Model{ID: tokenResult.Token[0].Model.ID},
		ExpiredAt: time.Now().Add(30 * time.Minute), UID: userResult.User[0].ID, Status: 1, AccessToken: accessToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		tmp := []token.Token{{AccessToken: accessToken}}
		c.JSON(http.StatusOK, &token.Result{Status: true, Token: tmp})
	}
}

func Delete(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := dbToken.Get(token.UserTokenAndAccessToken, &token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: result.Err.Error()})
		return
	}
	if err := dbToken.Delete(&token.Token{Model: gorm.Model{ID: result.Token[0].ID}}); err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Status: true})
}