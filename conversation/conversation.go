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
	Keys  []string
	times = 0
)

func init() {
	fileBytes, err := os.ReadFile("key.txt")
	if err != nil {
		fmt.Println("Failed to Read key from key.txt")
		fmt.Println(err.Error())
		os.Exit(0)
	}
	readLines := strings.Split(string(fileBytes), "\n")
	if len(readLines) == 0 {
		panic("No valid keys in key.txt")
	}
	Keys = make([]string, len(readLines))
	keyNums := 0
	for i := range readLines {
		if len(readLines[i]) > 2 {
			readLines[i] = strings.Trim(readLines[i], " ")
			Keys[keyNums] = readLines[i]
			fmt.Printf("Loaded Key: %s\n", Keys[keyNums])
			keyNums++
		}
	}
	if keyNums == 0 {
		panic("No Valid Keys! Make sure you have put keys in key.txt")
	}
	if keyNums != cap(Keys) {
		Keys = append(Keys[:keyNums])
	}
}

type Conversation struct {
	Prompt          string     // prompt
	SentenceList    *list.List // 对话表, 偶数是人类
	AIAnswered      bool       // AI是否完成回复
	LastModify      int64      // 上次回复的时间戳
	RequestSettings *RequestSettings
	AIName          string
	HumanName       string
}

// CStorage :Conversation存储形式
type CStorage struct {
	Id           string
	Prompt       string   // prompt
	SentenceList []string // 对话表, 偶数是人类
}
type RequestSettings struct {
	Model            string   `json:"model"`
	MaxTokens        int      `json:"max_tokens"`
	Temperature      float32  `json:"temperature"`
	TopP             int      `json:"top_p"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
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

func SendDirectly(prompt string, settings *RequestSettings) (string, error) {
	// 构造请求头
	key := Key()
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", key)
	// 请求体, 转化为bytes
	request := ChatgptRequest{
		Prompt:          prompt,
		RequestSettings: *settings,
	}
	jsonString, err := json.Marshal(&request)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 发送请求
	result, err := util.PostHeader("https://api.openai.com/v1/completions", jsonString, headers)
	if err != nil {
		if err.Error() == "Post \"https://api.openai.com/v1/completions\": net/http: invalid header field value for \"Authorization\"" && flgs.AutoRemoveErrorKeys {
			findAndRemoveKey(key)
		}
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
		err = ChatgptError.Err(jsonObject.S("error", "message").Data().(string))
		fmt.Println(err)
		switch err.(type) {
		case ChatgptError.ExceededQuotaException:
			if flgs.AutoRemoveErrorKeys {
				findAndRemoveKey(key)
			}
		}
		return "", err
	}
	answer := answerData.(string)
	return answer, nil
}
func (conv *Conversation) GetAnswer(question string) (string, error) {
	if !conv.AIAnswered {
		return "", ChatgptError.ChatgptError{Msg: "AI is thinking"}
	}
	conv.AIAnswered = false
	conv.SentenceList.PushBack(question)
	key := Key()
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", key)
	prompt := conv.PlainText() + conv.AIName
	request := ChatgptRequest{
		Prompt:          prompt,
		RequestSettings: *conv.RequestSettings,
	}
	jsonString, err := json.Marshal(&request)
	defer func() {
		conv.AIAnswered = true
		conv.LastModify = time.Now().Unix()
	}()
	if err != nil {
		conv.SentenceList.Remove(conv.SentenceList.Back())
		fmt.Println(err)
		return "", err
	}
	result, err := util.PostHeader("https://api.openai.com/v1/completions", jsonString, headers)
	if err != nil {
		conv.SentenceList.Remove(conv.SentenceList.Back())
		if err.Error() == "Post \"https://api.openai.com/v1/completions\": net/http: invalid header field value for \"Authorization\"" && flgs.AutoRemoveErrorKeys {
			findAndRemoveKey(key)
		}
		fmt.Println(err)
		return "", err
	}
	jsonObject, err := gabs.ParseJSON([]byte(result))
	if err != nil {
		conv.SentenceList.Remove(conv.SentenceList.Back())
		fmt.Println(err)
		return "", err
	}
	answerData := jsonObject.S("choices", "0", "text").Data()
	if answerData == nil {
		conv.SentenceList.Remove(conv.SentenceList.Back())
		err = ChatgptError.Err(jsonObject.S("error", "message").Data().(string))
		fmt.Println(err.Error())
		switch err.(type) {
		case ChatgptError.ExceededQuotaException:
			if flgs.AutoRemoveErrorKeys {
				findAndRemoveKey(key)
			}
		}
		return "", err
	}
	answer := strings.Trim(answerData.(string), "\n")
	conv.SentenceList.PushBack(answer)
	return answer, nil
}

func CreateCustomConversation(prompt string, settings *RequestSettings, AIName string, HumanName string) *Conversation {
	return &Conversation{prompt, list.New(), true, time.Now().Unix(), settings, AIName, HumanName}
}

func CreateQuickConversation(prompt string, sentenceList *list.List) *Conversation {
	return &Conversation{prompt, sentenceList, true, time.Now().Unix(), QuickChatSettings, "\nAI: ", "\nHuman: "}
}

func (conv *Conversation) PlainText() string {
	builder := new(strings.Builder)
	builder.WriteString(conv.Prompt)
	builder.WriteString("\n")
	index := 0
	for element := conv.SentenceList.Front(); element != nil; element = element.Next() {
		if index%2 == 0 {
			builder.WriteString(conv.HumanName)
		} else {
			builder.WriteString(conv.AIName)
		}
		builder.WriteString(element.Value.(string))
		index++
	}
	return builder.String()
}

func (conv *Conversation) ToCStorage(id string) *CStorage {
	sentences := make([]string, conv.SentenceList.Len())
	i := 0
	for element := conv.SentenceList.Front(); element != nil; element = element.Next() {
		sentences[i] = element.Value.(string)
		i++
	}
	return &CStorage{id, conv.Prompt, sentences}
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
	if len(Keys) <= 1 {
		fmt.Println("No Valid Keys, server stopped")
		os.Exit(0)
	}
	for i := range Keys {
		if Keys[i] == key {
			fmt.Printf("Key is not available and has been removed: %s\n", Keys[i])
			Keys = append(Keys[:i], Keys[i+1:]...)
			return
		}
	}
}
