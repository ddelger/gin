package persistence

import "github.com/ddelger/gin/model"

type Manager interface {
	UserManager

	Close()
}

type UserManager interface {
	GetUserById(int64, *model.User) error
	GetUserByEmail(string, *model.User) error
	GetUserByToken(string, *model.User) error

	ExistsUser(string) bool

	DeleteUser(*model.User) error
	UpdateUser(*model.User) error
	CreateUser(*model.User) error
}
