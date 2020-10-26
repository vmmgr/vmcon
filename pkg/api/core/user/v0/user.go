package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	auth "github.com/vmmgr/controller/pkg/api/core/auth/v0"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/core/tool/config"
	"github.com/vmmgr/controller/pkg/api/core/tool/gen"
	"github.com/vmmgr/controller/pkg/api/core/tool/hash"
	"github.com/vmmgr/controller/pkg/api/core/tool/mail"
	"github.com/vmmgr/controller/pkg/api/core/tool/notification"
	toolToken "github.com/vmmgr/controller/pkg/api/core/tool/token"
	"github.com/vmmgr/controller/pkg/api/core/user"
	dbUser "github.com/vmmgr/controller/pkg/api/store/user/v0"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Add(c *gin.Context) {
	var input, data user.User
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong email address")})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong name")})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}

	mailToken, _ := toolToken.Generate(4)

	pass := ""

	// 新規ユーザ
	if input.GroupID == 0 { //new user
		if input.Pass == "" {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong pass")})
			return
		}
		data = user.User{GroupID: 0, Name: input.Name, Email: input.Email, Pass: input.Pass, Status: 0, Level: 1,
			MailVerify: false, MailToken: mailToken}

		// グループ所属ユーザの登録
	} else {
		if input.Level == 0 || input.Level > 5 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong user level")})
			return
		}
		authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
		if authResult.Err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
			return
		}
		if authResult.User.GroupID != input.GroupID && authResult.User.GroupID > 0 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "GroupID mismatch"})
			return
		}

		pass = gen.GenerateUUID()
		log.Println("Email: " + input.Email)
		log.Println("tmp_Pass: " + pass)

		data = user.User{GroupID: input.GroupID, Name: input.Name, Email: input.Email, Status: 0, Level: input.Level,
			Pass: strings.ToLower(hash.Generate(pass)), MailVerify: false, MailToken: mailToken}
	}

	//check exist for database
	if err := dbUser.Create(&data); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		attachment := slack.Attachment{}
		attachment.AddField(slack.Field{Title: "E-Mail", Value: input.Email}).
			AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(input.GroupID))}).
			AddField(slack.Field{Title: "Name", Value: input.Name})

		notification.SendSlack(notification.Slack{Attachment: attachment, Channel: "user", Status: true})

		if pass == "" {
			mail.SendMail(mail.Mail{
				ToMail:  data.Email,
				Subject: "本人確認のメールにつきまして",
				Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
					config.Conf.Controller.User.Url + "/api/v1/user/verify/" + mailToken + "\n" +
					"本人確認が完了次第、ログイン可能になります。\n",
			})
		} else {
			mail.SendMail(mail.Mail{
				ToMail:  data.Email,
				Subject: "本人確認メールにつきまして",
				Content: " " + input.Name + "様\n\n" + "以下のリンクから本人確認を完了してください。\n" +
					config.Conf.Controller.User.Url + "/api/v1/user/verify/" + mailToken + "\n" +
					"本人確認が完了次第、ログイン可能になります。\n" + "仮パスワード: " + pass,
			})
		}

		c.JSON(http.StatusOK, user.Result{Status: true})
	}
}

func MailVerify(c *gin.Context) {
	token := c.Param("token")

	result := dbUser.Get(user.MailToken, &user.User{MailToken: token})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: result.Err.Error() + "| we can't find token data"})
		return
	}

	if result.User[0].MailVerify {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("This email has already been checked")})
		return
	}
	if result.User[0].Status >= 100 {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("error: user status")})
		return
	}

	if err := dbUser.Update(user.UpdateVerifyMail, &user.User{Model: gorm.Model{ID: result.User[0].ID},
		MailVerify: true}); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, &user.Result{Status: true})
	}
}

func Update(c *gin.Context) {
	var input user.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: fmt.Sprintf("id error")})
		return
	}
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
		return
	}

	if !authResult.User.MailVerify {
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: "not verify for user mail"})
		return
	}

	var u, serverData user.User

	if authResult.User.ID == uint(id) || id == 0 {
		serverData = authResult.User
		u.Model.ID = authResult.User.ID
		u.Status = authResult.User.Status
	} else {
		if authResult.User.GroupID == 0 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: Group ID = 0"})
			return
		}
		if authResult.User.Level > 1 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: failed user level"})
			return
		}
		userResult := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: uint(id)}})
		if userResult.Err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: userResult.Err.Error()})
			return
		}
		if userResult.User[0].GroupID != authResult.User.GroupID {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("failed group authentication")})
			return
		}
		serverData = userResult.User[0]
		u.Model.ID = uint(id)
		u.Status = userResult.User[0].Status
	}

	u, err = replaceUser(serverData, input, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbUser.Update(user.UpdateInfo, &u); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, user.Result{Status: true})
	}
}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	authResult.User.Pass = ""
	authResult.User.MailToken = ""
	if authResult.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
	} else {
		c.JSON(http.StatusOK, user.ResultOne{Status: true, User: authResult.User})
	}
}

func GetGroup(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	result := dbUser.Get(user.GID, &user.User{GroupID: authResult.Group.ID})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: result.Err.Error()})
		return
	}

	var data []user.User

	for _, tmp := range result.User {
		tmp.Pass = ""
		tmp.MailToken = ""
		if 0 < tmp.Status && tmp.Status < 100 {
			data = append(data, tmp)
		}
	}
	c.JSON(http.StatusOK, user.Result{Status: true, User: data})
}