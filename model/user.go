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
	Name     string `json:"name"`
	Token    string `json:"token"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Pricing  int64  `json:"pricing"`
	Created  int64  `json:"created"`
}

func (u *User) String() string {
	return u.Name
}

func NewUser(name, email, password string) *User {
	return &User{
		Token:    bson.NewObjectId().Hex(),
		Name:     name,
		Email:    email,
		Password: password,
		Pricing:  PricingPlanDefault,
		Created:  time.Now().Unix(),
	}
}
