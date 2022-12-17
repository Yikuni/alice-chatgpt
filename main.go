package main

import (
	"alice-chatgpt/conversation"
	"alice-chatgpt/util"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

var (
	idLength        = 8
	token           string
	conversationMap = make(map[string]*conversation.Conversation, 20)
)

func main() {
	flag.StringVar(&token, "t", "alice", "verify token")
	app := gin.Default()
	app.POST("/chatgpt/create", create)
	app.POST("/chatgpt/chat", chat)
	app.POST("/chatgpt/context", context)
	app.POST("/chatgpt/finish", finish)
	err := app.Run(":7777")
	if err != nil {
		fmt.Println(err)
	}
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

/*
*
创建会话
请求体: prompt plain text
*/
func create(c *gin.Context) {
	if !verify(c) {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	body := string(bodyBytes)
	runes := util.RandStringRunes(idLength)
	// 保证没有重复runes
	for ; conversationMap[runes] != nil; runes = util.RandStringRunes(idLength) {
	}
	if body == "" {
		conversationMap[runes] = conversation.CreateDefaultConversation()
	} else {
		conversationMap[runes] = conversation.CreateConversation(body)
	}
	c.String(200, runes)
}

/*
*
获取一个会话的上文
*/
func context(c *gin.Context) {
	if !verify(c) {
		return
	}

	conv := getConversation(c)
	if conv == nil {
		return
	}
	c.String(200, conv.PlainText())
}

func chat(c *gin.Context) {
	if !verify(c) {
		return
	}
	conv := getConversation(c)
	if conv == nil {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	body := string(bodyBytes)
	answer, err := conv.GetAnswer(body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, answer)
}

func finish(c *gin.Context) {
	if !verify(c) {
		return
	}
	conv := getConversation(c)
	if conv == nil {
		return
	}
	delete(conversationMap, c.GetHeader("conversation"))
	c.String(200, "Conversation finished")
}
