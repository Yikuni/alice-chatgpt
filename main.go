package main

import (
	"alice-chatgpt/conversation"
	"flag"
	"github.com/gin-gonic/gin"
)

var token string

var conversationMap = make(map[string]*conversation.Conversation, 20)

func main() {
	flag.StringVar(&token, "t", "alice", "verify token")
	app := gin.Default()
	app.POST("/chatgpt/create", create)
	app.POST("/chatgpt/createWithContext", createWithContext)
	app.POST("/chatgpt/chat", chat)
	app.POST("/chatgpt/context", context)
	app.POST("/chatgpt/finish", finish)
}

func verify(c *gin.Context) bool {
	if token != c.GetHeader("token") {
		c.String(505, "Auth failed")
		return false
	} else {
		return true
	}
}

func getConversation(c *gin.Context) *conversation.Conversation {
	id := c.GetHeader("conversation")
	if id == "" {
		c.String(500, "Conversation id is nil")
		return nil
	}
	conv := conversationMap[id]
	if conv == nil {
		c.String(500, "failed to find Conversation with provided id: %s ", id)
		return nil
	}
	return conv
}

// 创建会话
func create(c *gin.Context) {
	verify(c)
}

func context(c *gin.Context) {
	verify(c)

}

func createWithContext(c *gin.Context) {
	verify(c)

}

func chat(c *gin.Context) {
	verify(c)

}

func finish(c *gin.Context) {
	verify(c)
}
