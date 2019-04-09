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

	if len(b.Username) == 0 || len(b.Password) == 0 {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid username [%s] or password.", b.Username))
		return
	}

	u := &model.User{}
	if err := m.GetUserByName(b.Username, u); err != nil {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid username [%s] or password.", b.Username))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(b.Password)); err != nil {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid username [%s] or password.", b.Username))
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

	if len(b.Username) == 0 || len(b.Password) == 0 {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("Invalid username [%s] or password.", b.Username))
		return
	}

	if m.ExistsUser(b.Username) {
		AppendResponseError(c, http.StatusOK, fmt.Errorf("User [%s] already exists.", b.Username))
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(b.Password), bcrypt.MinCost)
	if err != nil {
		glog.Errorf("Unable to encrypt password. %s", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	u := model.NewUser(b.Username, string(h))
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
			Id:        u.Username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
		User: &model.User{
			Token:    u.Token,
			Pricing:  u.Pricing,
			Created:  u.Created,
			Username: u.Username,
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
