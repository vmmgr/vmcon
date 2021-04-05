package user

import "github.com/vmmgr/controller/pkg/api/core"

const (
	ID               = 0
	GID              = 1
	Name             = 2
	Email            = 3
	MailToken        = 4
	UpdateVerifyMail = 100
	UpdateGroupID    = 101
	UpdateInfo       = 102
	UpdateLevel      = 106
	UpdateAll        = 110
)

type ResultOne struct {
	Status bool      `json:"status"`
	Error  string    `json:"error"`
	User   core.User `json:"data"`
}

type Result struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	User   []core.User `json:"data"`
}

type ResultDatabase struct {
	Err  error
	User []core.User
}
