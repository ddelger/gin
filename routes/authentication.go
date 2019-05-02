package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ddelger/gin/middleware"
	"github.com/ddelger/gin/model"
	"github.com/ddelger/gin/persistence"

	"github.com/ddelger/glog"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	User  = "user"
	Token = "token"
)

func Login(c *gin.Context) {
	m := (c.MustGet(middleware.MiddlewarePersistence)).(persistence.Manager)

	b := &model.User{}
	if err := c.BindJSON(b); err != nil {
		AppendResponseError(c, http.StatusBadRequest, err)
		return
	}

	if len(b.Email) == 0 || len(b.Password) == 0 {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid email [%s] or password.", b.Email))
		return
	}

	u := &model.User{}
	if err := m.GetUserByEmail(b.Email, u); err != nil {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid email [%s] or password.", b.Email))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(b.Password)); err != nil {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid email [%s] or password.", b.Email))
		return
	}

	IssueToken(c, u)
}

func Register(c *gin.Context) {
	m := (c.MustGet(middleware.MiddlewarePersistence)).(persistence.Manager)

	b := &model.User{}
	if err := c.BindJSON(b); err != nil {
		AppendResponseError(c, http.StatusBadRequest, err)
		return
	}

	if len(b.Name) == 0 || len(b.Email) == 0 || len(b.Password) == 0 {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid email [%s] or password.", b.Email))
		return
	}

	if m.ExistsUser(b.Email) {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Email [%s] already exists.", b.Email))
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(b.Password), bcrypt.MinCost)
	if err != nil {
		glog.Errorf("Unable to encrypt password. %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	u := model.NewUser(b.Name, b.Email, string(h))
	if err := m.CreateUser(u); err != nil {
		glog.Errorf("Unable to create user. %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	IssueToken(c, u)
}

func IssueToken(c *gin.Context, u *model.User) {
	k := (c.MustGet(middleware.MiddlewareKeys)).(*model.Keys)

	t := jwt.New(jwt.SigningMethodRS512)
	t.Claims = &model.Claims{
		StandardClaims: &jwt.StandardClaims{
			Id:        u.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		User: &model.User{
			Name:    u.Name,
			Email:   u.Email,
			Token:   u.Token,
			Pricing: u.Pricing,
			Created: u.Created,
		},
	}

	s, err := t.SignedString(k.PrivateKey)
	if err != nil {
		glog.Errorf("Unable to generate token. %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{Token: s})
}
