package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core/token"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(t *token.Token) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Create(t).Error
}

func Delete(t *token.Token) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(t).Error
}

func DeleteAll() error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Exec("DELETE FROM tokens").Error
}

func Update(base int, t *token.Token) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	if token.AddToken == base {
		err = db.Model(&token.Token{Model: gorm.Model{ID: t.ID}}).Update(token.Token{Model: gorm.Model{},
			ExpiredAt: t.ExpiredAt, UID: t.UID, Status: t.Status, AccessToken: t.AccessToken}).Error
	} else if token.UpdateToken == base {
		err = db.Model(&token.Token{Model: gorm.Model{ID: t.ID}}).Update("expired_at", t.ExpiredAt).Error
	} else if token.UpdateAll == base {
		err = db.Model(&token.Token{Model: gorm.Model{ID: t.ID}}).Update(token.Token{
			ExpiredAt: t.ExpiredAt, UID: t.UID, Status: t.Status, UserToken: t.UserToken, TmpToken: t.TmpToken,
			AccessToken: t.AccessToken, Debug: t.Debug}).Error
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n %s", time.Now(), err)
	}
	return err
}

// value of base can reference from api/core/user/interface.go
func Get(base int, input *token.Token) token.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return token.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var tokenStruct []token.Token

	if base == token.UserToken {
		err = db.Where("user_token = ? AND admin = ? AND expired_at > ?",
			input.UserToken, false, time.Now()).Find(&tokenStruct).Error
	} else if base == token.UserTokenAndAccessToken {
		err = db.Where("user_token = ? AND access_token = ? AND admin = ? AND expired_at > ?",
			input.UserToken, input.AccessToken, false, time.Now()).Find(&tokenStruct).Error
	} else if base == token.AdminToken {
		err = db.Where("access_token = ? AND admin = ? AND expired_at > ?",
			input.AccessToken, true, time.Now()).Find(&tokenStruct).Error
	} else if base == token.ExpiredTime {
		err = db.Where("expired_at < ? ", time.Now()).Find(&tokenStruct).Error
	} else {
		log.Println("base select error")
		return token.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return token.ResultDatabase{Token: tokenStruct, Err: err}
}

func GetAll() token.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return token.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var tokens []token.Token
	err = db.Find(&tokens).Error
	return token.ResultDatabase{Token: tokens, Err: err}
}