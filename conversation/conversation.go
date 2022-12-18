package conversation

import (
	"alice-chatgpt/util"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"os"
	"strings"
	"time"
)

var Key string

func init() {
	fileBytes, err := os.ReadFile("key.txt")
	if err != nil {
		fmt.Println("Failed to Read key from key.txt")
		fmt.Println(err.Error())
		os.Exit(0)
	}
	Key = string(fileBytes)

}

type Conversation struct {
	Prompt       string     // prompt
	SentenceList *list.List // 对话表, 偶数是人类
	AIAnswered   bool       // AI是否完成回复
	LastModify   int64      // 上次回复的时间戳
}

type ChatgptRequest struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
	PresencePenalty  float32 `json:"presence_penalty"`
}

func (conversation *Conversation) GetAnswer(question string) (string, error) {
	conversation.AIAnswered = false
	conversation.SentenceList.PushBack(question)
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + Key
	request := ChatgptRequest{
		Model:            "text-davinci-003",
		Prompt:           conversation.PlainText() + "\nAI: ",
		MaxTokens:        10000,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Temperature:      0.9}
	jsonString, err := json.Marshal(&request)
	defer func() {
		conversation.AIAnswered = true
		conversation.LastModify = time.Now().Unix()
	}()
	if err != nil {
		conversation.SentenceList.Remove(conversation.SentenceList.Back())
		fmt.Println(err)
		return "", err
	}
	result, err := util.PostHeader("https://api.openai.com/v1/completions", jsonString, headers)
	if err != nil {
		conversation.SentenceList.Remove(conversation.SentenceList.Back())
		fmt.Println(err)
		return "", err
	}
	jsonObject, err := gabs.ParseJSON([]byte(result))
	if err != nil {
		conversation.SentenceList.Remove(conversation.SentenceList.Back())
		fmt.Println(err)
		return "", err
	}
	fmt.Println(jsonObject.Data())
	answers := jsonObject.S("choices", "*", "text")
	fmt.Println(answers.String())
	answer := strings.Builder{}
	children := answers.Children()
	for i := range children {
		answer.WriteString(children[i].String())
		fmt.Printf("i = %d, answer = %s\n", i, children[i].String())
	}
	answerString := answer.String()
	fmt.Printf("result: %s", answerString)
	conversation.SentenceList.PushBack(answerString)
	return answerString, nil
}

func CreateConversation(prompt string) *Conversation {
	return &Conversation{prompt, list.New(), true, time.Now().Unix()}
}

func CreateDefaultConversation() *Conversation {
	return CreateConversation("The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n")
}

func (conversation *Conversation) PlainText() string {
	builder := new(strings.Builder)
	builder.WriteString(conversation.Prompt)
	builder.WriteString("\n")
	index := 0
	for element := conversation.SentenceList.Front(); element != nil; element = element.Next() {
		if index%2 == 0 {
			builder.WriteString("\nHuman: ")
		} else {
			builder.WriteString("\nAI: ")
		}
		builder.WriteString(element.Value.(string))
		index++
	}
	return builder.String()
}
