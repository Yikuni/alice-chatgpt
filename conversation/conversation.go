package conversation

import (
	"alice-chatgpt/util"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"os"
	"strings"
)

var Key string

func init() {
	fileBytes, err := os.ReadFile("key.txt")
	if err != nil {
		fmt.Println("Failed to Read key: key.txt")
		fmt.Println(err.Error())
		return
	}
	Key = string(fileBytes)

}

type Conversation struct {
	Prompt       string     // prompt
	SentenceList *list.List // 对话表, 偶数是人类
	AIAnswered   bool       // AI是否完成回复
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
		MaxTokens:        150,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Temperature:      0.9}
	jsonString, err := json.Marshal(&request)
	defer func() { conversation.AIAnswered = true }()
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
	answer := jsonObject.S("choices", "*", "text").Data()
	conversation.SentenceList.PushBack(answer.(string))
	return answer.(string), nil
}

func CreateConversation(prompt string) *Conversation {
	return &Conversation{prompt, list.New(), true}
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
