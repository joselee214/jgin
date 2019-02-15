package jgin

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	storesession "github.com/gorilla/sessions"
	sessions "github.com/tommy351/gin-sessions"
)

func SetSession(ctx *gin.Context, k string, o interface{}) {
	session := sessions.Get(ctx)
	session.Set(k, o)
	session.Save()
}

func GetSession(ctx *gin.Context, k string) interface{} {
	session := sessions.Get(ctx)
	return session.Get(k)
}

func ClearAllSession(ctx *gin.Context) {
	session := sessions.Get(ctx)
	session.Clear()
	session.Save()
	return
}

//修改 session 以便传入MaxAge
func Typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func NewCookieStore(maxAge int, keyPairs ...[]byte) sessions.CookieStore {

	cs := &storesession.CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &storesession.Options{
			Path:   "/",
			MaxAge: maxAge,
		},
	}
	cs.MaxAge(cs.Options.MaxAge)

	fmt.Println("==================", Typeof(cs), cs)

	return &cookieStore{cs}
}

type cookieStore struct {
	*storesession.CookieStore
}

func (c *cookieStore) Options(options sessions.Options) {
	c.CookieStore.Options = &storesession.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
