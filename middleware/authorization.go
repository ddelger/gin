package middleware

import (
	"io/ioutil"
	"net/http"

	"github.com/ddelger/gin/model"
	"github.com/ddelger/gin/persistence"

	"github.com/ddelger/glog"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

const (
	MiddlewareKeys = "MiddlewareKeys"
	MiddlewareUser = "MiddlewareUser"
)

func Keys(keys *model.Keys) gin.HandlerFunc {
	f, err := ioutil.ReadFile(keys.PrivateKeyFile)
	if err != nil {
		glog.Fatalf("Unable to read private key [%s]. %s", keys.PrivateKeyFile, err)
	}
	prv, err := jwt.ParseRSAPrivateKeyFromPEM(f)
	if err != nil {
		glog.Fatalf("Unable to parse private key [%s]. %s", keys.PrivateKeyFile, err)
	}

	f, err = ioutil.ReadFile(keys.PublicKeyFile)
	if err != nil {
		glog.Fatalf("Unable to read public key [%s]. %s", keys.PublicKeyFile, err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(f)
	if err != nil {
		glog.Fatalf("Unable to parse public key [%s]. %s", keys.PublicKeyFile, err)
	}

	return func(c *gin.Context) {
		c.Set(MiddlewareKeys, &model.Keys{PublicKey: pub, PrivateKey: prv})
		c.Next()
	}
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		k := (c.MustGet(MiddlewareKeys)).(*model.Keys)
		m := (c.MustGet(MiddlewarePersistence)).(persistence.Manager)

		t, err := request.ParseFromRequestWithClaims(c.Request, request.OAuth2Extractor, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return k.PublicKey, nil
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		u := &model.User{}
		if err := m.GetUserByName(t.Claims.(*model.Claims).Id, u); err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set(MiddlewareUser, u)
		c.Next()
	}
}
