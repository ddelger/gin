package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	PricingPlanDefault = 0
)

type User struct {
	Id       int64  `bson:"_id" json:"id"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Pricing  int64  `json:"pricing"`
	Created  int64  `json:"created"`
}

func (u *User) String() string {
	return u.Username
}

func NewUser(username, password string) *User {
	return &User{
		Token:    bson.NewObjectId().Hex(),
		Username: username,
		Password: password,
		Pricing:  PricingPlanDefault,
		Created:  time.Now().Unix(),
	}
}
