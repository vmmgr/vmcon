package v0

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/pkg/api/core"
	"github.com/vmmgr/controller/pkg/api/core/support/ticket"
	"github.com/vmmgr/controller/pkg/api/store"
	"log"
	"time"
)

func Create(support *core.Ticket) (*core.Ticket, error) {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return support, fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	err = db.Create(&support).Error
	return support, err
}

func Delete(support *core.Ticket) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	return db.Delete(support).Error
}

func Update(base int, t core.Ticket) error {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())
	}
	defer db.Close()

	var result *gorm.DB

	//#4 Issue(解決済み）
	if ticket.UpdateAll == base {
		result = db.Model(&core.Ticket{Model: gorm.Model{ID: t.ID}}).Update(&core.Ticket{Title: t.Title,
			GroupID: t.GroupID,
			UserID:  t.UserID,
			Solved:  t.Solved,
		})
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n", time.Now())
	}
	return result.Error
}

func Get(base int, data *core.Ticket) ticket.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var ticketStruct []core.Ticket

	if base == ticket.ID { //ID
		err = db.Preload("User").
			Preload("Group").
			Preload("Chat").
			Preload("Chat.User").
			First(&ticketStruct, data.ID).Error
	} else if base == ticket.GID { //GroupID
		err = db.Where("group_id = ?", data.GroupID).
			Preload("User").
			Preload("Group").
			Preload("Chat").
			Preload("Chat.User").
			Find(&ticketStruct).Error
	} else if base == ticket.UID { //UserID
		err = db.Where("user_id = ?", data.UserID).Find(&ticketStruct).Error
	} else {
		log.Println("base select error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: base select\n", time.Now())}
	}
	return ticket.ResultDatabase{Tickets: ticketStruct, Err: err}
}

func GetAll() ticket.ResultDatabase {
	db, err := store.ConnectDB()
	if err != nil {
		log.Println("database connection error")
		return ticket.ResultDatabase{Err: fmt.Errorf("(%s)error: %s\n", time.Now(), err.Error())}
	}
	defer db.Close()

	var tickets []core.Ticket
	err = db.Preload("User").
		Preload("Group").
		Preload("Chat").
		Preload("Chat.User").
		Find(&tickets).Error
	return ticket.ResultDatabase{Tickets: tickets, Err: err}
}
