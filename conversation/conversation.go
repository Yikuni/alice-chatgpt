package conversation

import (
	"alice-chatgpt/ChatgptError"
	flgs "alice-chatgpt/flags"
	"alice-chatgpt/util"
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"os"
	"strings"
	"time"
)

var (
	Keys            []string
	times           = 0
	DefaultSettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        500,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Temperature:      0.9,
	}
	QuickChatSettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        3500,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
		Temperature:      0.9,
	}
	SummarySettings = RequestSettings{
		Model:            "text-davinci-003",
		MaxTokens:        2048,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  1,
		Temperature:      0.7,
	}
)

func init() {
	fileBytes, err := os.ReadFile("key.txt")
	if err != nil {
		fmt.Println("Failed to Read key from key.txt")
		fmt.Println(err.Error())
		os.Exit(0)
	}
	Keys = strings.Split(string(fileBytes), "\n")
	if len(Keys) == 0 {
		panic("No valid keys in key.txt")
	}
	for i := range Keys {
		if i != len(Keys)-1 {
			last := Keys[i][len(Keys[i])-1]
			if last == '\n' {
				Keys[i] = Keys[i][:len(Keys[i])-1]
			}
		}
	}
}

type Conversation struct {
	Prompt       string     // prompt
	SentenceList *list.List // 对话表, 偶数是人类
	AIAnswered   bool       // AI是否完成回复
	LastModify   int64      // 上次回复的时间戳
}

// CStorage :Conversation存储形式
type CStorage struct {
	Id           string
	Prompt       string   // prompt
	SentenceList []string // 对话表, 偶数是人类
}
type RequestSettings struct {
	Model            string  `json:"model"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float32 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty float32 `json:"frequency_penalty"`
	PresencePenalty  float32 `json:"presence_penalty"`
}

type ChatgptRequest struct {
	RequestSettings
	Prompt string `json:"prompt"`
}

func Key() string {
	key := Keys[times%len(Keys)]
	times++
	return key
}

func SendDirectly(prompt string, settings RequestSettings) (string, error) {
	// 构造请求头
	key := Key()
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + key
	// 请求体, 转化为bytes
	request := ChatgptRequest{
		Prompt:          prompt,
		RequestSettings: settings,
	}
	jsonString, err := json.Marshal(&request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 发送请求
	result, err := util.PostHeader("https://api.openai.com/v1/completions", jsonString, headers)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 处理请求, 获得文字结果
	jsonObject, err := gabs.ParseJSON([]byte(result))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(jsonObject.String())
	answerData := jsonObject.S("choices", "0", "text").Data()
	if answerData == nil {
		fmt.Println(err)
		switch err.(type) {
		case ChatgptError.ExceededQuotaException:
			if flgs.AutoRemoveErrorKeys {
				findAndRemoveKey(key)
			}
		}

		return "", ChatgptError.Err(jsonObject.S("error", "message").Data().(string))
	}
	answer := answerData.(string)
	return answer, nil
}
func (conversation *Conversation) GetAnswer(question string, settings RequestSettings) (string, error) {
	conversation.AIAnswered = false
	conversation.SentenceList.PushBack(question)
	key := Key()
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + key
	request := ChatgptRequest{
		Prompt:          conversation.PlainText() + "\nAI: ",
		RequestSettings: settings,
	}
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
	answerData := jsonObject.S("choices", "0", "text").Data()
	if answerData == nil {
		conversation.SentenceList.Remove(conversation.SentenceList.Back())
		fmt.Println(err)
		switch err.(type) {
		case ChatgptError.ExceededQuotaException:
			if flgs.AutoRemoveErrorKeys {
				findAndRemoveKey(key)
			}
		}

		return "", ChatgptError.Err(jsonObject.S("error", "message").Data().(string))
	}
	answer := answerData.(string)
	conversation.SentenceList.PushBack(answer)
	return answer, nil
}

func CreateConversation(prompt string) *Conversation {
	return &Conversation{prompt, list.New(), true, time.Now().Unix()}
}

func CreateQuickConversation(prompt string, sentenceList *list.List) *Conversation {
	return &Conversation{prompt, sentenceList, true, time.Now().Unix()}
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

func (conversation *Conversation) ToCStorage(id string) *CStorage {
	sentences := make([]string, conversation.SentenceList.Len())
	i := 0
	for element := conversation.SentenceList.Front(); element != nil; element = element.Next() {
		sentences[i] = element.Value.(string)
		i++
	}
	return &CStorage{id, conversation.Prompt, sentences}
}

// ToJsonBytes 转换为json byte
func (cStorage *CStorage) ToJsonBytes() ([]byte, error) {
	return json.Marshal(&cStorage)
}

// FromJsonBytes 通过jsonBytes获得CStorage对象
func FromJsonBytes(marshal []byte) (*CStorage, error) {
	jsonObj, err := gabs.ParseJSON(marshal)
	if err != nil {
		return nil, err
	}
	children := jsonObj.S("SentenceList").Children()
	sentences := make([]string, len(children))
	for i, sentence := range children {
		sentences[i] = sentence.Data().(string)
	}
	return &CStorage{jsonObj.S("Id").Data().(string), jsonObj.S("Prompt").Data().(string), sentences}, nil
}

func findAndRemoveKey(key string) {
	for i := range Keys {
		if Keys[i] == key {
			Keys = append(Keys[:i], Keys[i+1:]...)
			return
		}
	}
}
