package auth

import (
	global "alice-chatgpt/global"
	"alice-chatgpt/util"
	"container/list"
	"crypto/sha256"
	"github.com/gin-gonic/gin"
)

type Auth interface {
	Verify(ctx *gin.Context) bool
}

// SimpleAuth simple
type SimpleAuth struct {
	Auth
}

// NoneAuth none
type NoneAuth struct {
	Auth
}

type NormalAuth struct {
	uuidList *list.List
	maxSize  int
	Auth
}

func NewNormalAuth(maxSize int) *NormalAuth {
	return &NormalAuth{maxSize: maxSize, uuidList: list.New()}
}

func (auth *SimpleAuth) Verify(ctx *gin.Context) bool {
	return ctx.GetHeader("token") == global.Token
}

func (auth *NoneAuth) Verify(ctx *gin.Context) bool {
	return true
}

func (auth *NormalAuth) Verify(ctx *gin.Context) bool {
	uuid := ctx.GetHeader("uuid")
	if len(uuid) > 50 {
		return false
	}
	if uuid == "" {
		return false
	}
	if auth.uuidList.Len() > auth.maxSize {
		auth.uuidList = list.New()
	} else {
		for i := auth.uuidList.Front(); i != nil; i = i.Next() {
			if uuid == i.Value {
				return false
			}
		}
		auth.uuidList.PushFront(uuid)
	}
	token := ctx.GetHeader("token")
	trueToken := sha256.Sum256([]byte(global.Token + uuid))
	return util.HexBuffToString(trueToken[:]) == token
}
