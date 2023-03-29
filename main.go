package main

import (
	"alice-chatgpt/auth"
	"alice-chatgpt/conversation"
	"alice-chatgpt/dao"
	"alice-chatgpt/global"
	"alice-chatgpt/util"
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

var (
	idLength = 8

	conversationMap = make(map[string]conversation.Conversation, 20)
	p               string
	db              string
	daoInstance     dao.Dao
	qps             = 0
	authInstance    auth.Auth
)

func main() {
	flag.StringVar(&global.Token, "t", "alice", "verify token")
	flag.StringVar(&p, "p", "7777", "port")
	flag.StringVar(&db, "db", "badger", "database url; use badger if undefined, now support only badger")
	flag.BoolVar(&global.AutoRemoveErrorKeys, "a", true, "auto remove key when key is above the quota")
	flag.IntVar(&global.LimitPerMin, "l", 500, "limit usage per minute")
	flag.StringVar(&global.AuthType, "auth", "simple", "auth type: none, simple, normal. simple as default")
	flag.StringVar(&global.Proxy, "proxy", "", "http proxy")
	flag.Parse()
	setAuthInstance()
	app := gin.Default()
	app.POST("/chatgpt/create", create)
	app.POST("/chatgpt/createTurbo", createTurbo)
	app.POST("/chatgpt/createGPT4", createGPT4)
	app.POST("/chatgpt/createRolePlay", createRolePlay)
	app.POST("/chatgpt/chat", chat)
	app.POST("/chatgpt/rollback", rollbackConversation)
	app.POST("/chatgpt/context", context)
	app.POST("/chatgpt/finish", finish)
	app.POST("/chatgpt/contextArray", contextArray)
	app.POST("/chatgpt/quickAnswer", quickAnswer)
	app.POST("/chatgpt/summary", summary)

	// 删除长时间没有使用的聊天
	go func() {
		for {
			time.Sleep(time.Minute * 30)
			now := time.Now().Unix()
			for k, v := range conversationMap {
				// 如果一定时间内没使用, 或者5分钟内仍然只有2句话
				duration := now - v.GetLastModify()
				if duration > 1200 || duration > 300 && v.GetSentenceList().Len() <= 2 {
					delete(conversationMap, k)
					fmt.Printf("GPT3Conversation with id: %s expired", k)
				}
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Minute)
			qps = 0
		}
	}()
	dbEnabled := false
	if db == "badger" {
		daoInstance = &dao.BadgerDao{}
	} else {
		// TODO Add mysql support
	}
	dbError := daoInstance.InitDatabase()
	if dbError != nil {
		fmt.Println(dbError.Error())
	} else {
		dbEnabled = true
	}
	if dbEnabled {
		app.GET("/chatgpt/share/:id", sharedContext)
		app.POST("/chatgpt/store", store)
		defer func() {
			err := daoInstance.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	} else {
		fmt.Println("Running without database!")
	}

	err := app.Run(":" + p)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func setAuthInstance() {
	switch global.AuthType {
	case "none":
		authInstance = &auth.NoneAuth{}
		fmt.Println("Auth set to none")
	case "normal":
		authInstance = auth.NewNormalAuth(50)
		fmt.Println("Auth set to normal")
	default:
		authInstance = &auth.SimpleAuth{}
		fmt.Println("Auth set to simple")
	}
}
func verify(c *gin.Context) bool {
	if authInstance.Verify(c) {
		return true
	} else {
		c.String(505, "Auth failed")
		return false
	}
}

func getConversation(c *gin.Context) *conversation.Conversation {
	id := c.GetHeader("conversation")
	if id == "" {
		c.String(500, "conversation id is nil")
		return nil
	}
	conv := conversationMap[id]
	if conv == nil {
		c.String(500, "failed to find conversation with provided id: %s ", id)
		return nil
	}
	return &conv
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

	settings := c.GetHeader("settings")
	var requestSettings *conversation.RequestSettings
	switch settings {
	case "":
		requestSettings = conversation.DefaultSettings
	case "FriendChat":
		requestSettings = conversation.FriendSettings
	default:
		c.String(500, "No such setting: "+settings)
		return
	}

	if body == "" && requestSettings == conversation.DefaultSettings {
		body = "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n"
	}
	conv := conversation.CreateCustomConversation(body, requestSettings, "\nAI: ", "\nHuman: ")
	if requestSettings == conversation.FriendSettings {
		conv.AIName = "\nFriend: "
		conv.HumanName = "\nYou: "
	}
	conversationMap[runes] = conv
	c.String(200, runes)
}

/*
*
创建会话
请求体: prompt plain text
*/
func createTurbo(c *gin.Context) {
	if !verify(c) {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if err != nil {
		c.String(500, err.Error())
		return
	}
	runes := util.RandStringRunes(idLength)
	// 保证没有重复runes
	for ; conversationMap[runes] != nil; runes = util.RandStringRunes(idLength) {
	}
	var (
		prompt    = "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."
		AIName    = "user"
		humanName = "assistant"
	)
	if len(bodyBytes) > 1 {
		prompt = string(bodyBytes)
	}
	conv := conversation.CreateTuborConversation(prompt, AIName, humanName)
	conversationMap[runes] = conv
	c.String(200, runes)
}

func createGPT4(c *gin.Context) {
	if !verify(c) {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	if err != nil {
		c.String(500, err.Error())
		return
	}
	runes := util.RandStringRunes(idLength)
	// 保证没有重复runes
	for ; conversationMap[runes] != nil; runes = util.RandStringRunes(idLength) {
	}
	var (
		prompt    = "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible."
		AIName    = "user"
		humanName = "assistant"
	)
	if len(bodyBytes) > 1 {
		prompt = string(bodyBytes)
	}
	conv := conversation.CreateGPT4Conversation(prompt, AIName, humanName)
	conversationMap[runes] = conv
	c.String(200, runes)
}

/*
*
create RolePlay conversation
*/
func createRolePlay(c *gin.Context) {
	if !verify(c) {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	/**
	{
		human_name: "Human",
		ai_name:	"AI",
		prompt:		"This is a conversation with an AI assistance..."
	}
	*/
	jsonObject, err := gabs.ParseJSON(bodyBytes)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	humanName := jsonObject.S("human_name").Data().(string)
	aiName := jsonObject.S("ai_name").Data().(string)
	prompt := jsonObject.S("prompt").Data().(string)
	runes := util.RandStringRunes(idLength)
	// 保证没有重复runes
	for ; conversationMap[runes] != nil; runes = util.RandStringRunes(idLength) {
	}
	requestSettings := conversation.RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        512,
		TopP:             1,
		FrequencyPenalty: 0.5,
		PresencePenalty:  0,
		Temperature:      0.5,
		Stop:             []string{humanName + ":"},
	}
	conv := conversation.CreateCustomConversation(prompt, &requestSettings, "\n"+aiName+": ", "\n"+humanName+": ")
	conversationMap[runes] = conv
	c.String(200, runes)
}

func rollbackConversation(c *gin.Context) {
	if !verify(c) {
		return
	}
	conv := getConversation(c)
	if conv == nil {
		return
	}
	err := conversation.Rollback(conv)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, "Rollback Succeeded")
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
	c.String(200, conversation.PlainText(*conv))
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
	answer, err := conversation.GetAnswer(*conv, body)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "" {
			errorMessage = "exceeded max tokens"
		}
		c.String(500, errorMessage)
		return
	}
	c.String(200, answer)
}

func quickAnswer(c *gin.Context) {
	if !verify(c) {
		return
	}
	if !limit(c) {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	jsonContainer, err := gabs.ParseJSON(bodyBytes)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	prompt := jsonContainer.S("prompt").Data().(string)
	question := jsonContainer.S("question").Data().(string)
	examples := list.New()
	for _, container := range jsonContainer.S("examples", "*").Children() {
		examples.PushBack(container.Data().(string))
	}
	conv := conversation.CreateQuickConversation(prompt, examples)
	answer, err := conversation.GetAnswer(conv, question)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "" {
			errorMessage = "exceeded max tokens"
		}
		c.String(500, errorMessage)
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
	c.String(200, "GPT3Conversation finished")
}

func contextArray(c *gin.Context) {
	if !verify(c) {
		return
	}
	conv := getConversation(c)
	if conv == nil {
		return
	}
	result, err := json.Marshal((*conv).GetSentenceList())
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, string(result))
}

// 保存会话, 用于分享
func store(c *gin.Context) {
	if !verify(c) {
		return
	}
	conv := getConversation(c)
	if conv == nil {
		return
	}
	cStorage := conversation.ToCStorage(*conv, c.GetHeader("conversation"))
	err := daoInstance.Save(cStorage)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.String(200, "保存成功")
}

// 获取分享的会话
func sharedContext(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(400, "Bad Request")
		return
	}
	cStorage := daoInstance.Search(id)
	if cStorage == nil {
		c.String(500, "Failed to find shared conversation with provided id")
		return
	}
	c.JSON(200, cStorage)
}

func summary(c *gin.Context) {
	if !verify(c) {
		return
	}
	if !limit(c) {
		return
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	article := string(bodyBytes)
	article += "\nTl;dr\n"
	answer, err := conversation.SendDirectly(article, conversation.SummarySettings)
	if err != nil {
		errorMessage := err.Error()
		if errorMessage == "" {
			errorMessage = "exceeded max tokens"
		}
		c.String(500, errorMessage)
		return
	}
	c.String(200, answer)
}

func limit(c *gin.Context) bool {
	if qps > global.LimitPerMin {
		c.String(502, "Server busy")
		return false
	}
	qps++
	return true
}
