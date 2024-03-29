package conversation

import (
	"alice-chatgpt/global"
	"alice-chatgpt/gpterror"
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

type Conversation interface {
	GetStreamFlag() bool
	GetSentenceList() *list.List
	GetPrompt() string
	GetHumanName() string
	GetAIName() string
	GetLastModify() int64
	SetLastModify(int64)
	GetAIAnswered() bool
	SetAIAnswered(bool)
	RequestBody() ([]byte, error)
	SolveResponse(jsonObject *gabs.Container) string
}

func GetAnswer(conv Conversation, question string) (string, error) {
	if question == "" {
		return "", gpterror.ChatgptError{Msg: "Question can not be empty"}
	}
	if !conv.GetAIAnswered() {
		return "", gpterror.ChatgptError{Msg: "AI is thinking"}
	}
	conv.SetAIAnswered(false)
	conv.GetSentenceList().PushBack(question)
	key := Key()
	headers := make(map[string]string, 2)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = fmt.Sprintf("%s %s", "Bearer", key)
	jsonString, err := conv.RequestBody()
	defer func() {
		conv.SetAIAnswered(true)
		conv.SetLastModify(time.Now().Unix())
	}()
	if err != nil {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		fmt.Println(err)
		return "", err
	}
	var _url string
	switch conv.(type) {
	case *GPT3Conversation:
		_url = "https://api.openai.com/v1/completions"
	case *TurboConversation:
		_url = "https://api.openai.com/v1/chat/completions"
	case *GPT4Conversation:
		_url = "https://api.openai.com/v1/chat/completions"
	}
	result, err := util.PostHeader(_url, jsonString, headers)
	if err != nil {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		if err.Error() == "Post \"https://api.openai.com/v1/completions\": net/http: invalid header field value for \"Authorization\"" && global.AutoRemoveErrorKeys {
			findAndRemoveKey(key)
		}
		fmt.Println(err)
		return "", err
	}
	jsonObject, err := gabs.ParseJSON([]byte(result))
	if err != nil {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		fmt.Println(err)
		return "", err
	}
	if jsonObject.Exists("error") {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		err = gpterror.Err(jsonObject.S("error", "message").Data().(string))
		fmt.Println(err.Error())
		switch err.(type) {
		case gpterror.ExceededQuotaException:
			if global.AutoRemoveErrorKeys {
				findAndRemoveKey(key)
			}
		}
		return "", err
	}
	answerData := conv.SolveResponse(jsonObject)
	answer := strings.Trim(answerData, "\n")
	conv.GetSentenceList().PushBack(answer)
	return answer, nil
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
		if err.Error() == "Post \"https://api.openai.com/v1/completions\": net/http: invalid header field value for \"Authorization\"" && global.AutoRemoveErrorKeys {
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
		err = gpterror.Err(jsonObject.S("error", "message").Data().(string))
		fmt.Println(err)
		switch err.(type) {
		case gpterror.ExceededQuotaException:
			if global.AutoRemoveErrorKeys {
				findAndRemoveKey(key)
			}
		}
		return "", err
	}
	answer := answerData.(string)
	return answer, nil
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
		fmt.Println("=================  No Valid Keys!  =================")
	}
	for i := range Keys {
		if Keys[i] == key {
			fmt.Printf("Key is not available and has been removed: %s\n", Keys[i])
			Keys = append(Keys[:i], Keys[i+1:]...)
			return
		}
	}
}

func ToCStorage(conv Conversation, id string) *CStorage {
	sentences := make([]string, conv.GetSentenceList().Len())
	i := 0
	for element := conv.GetSentenceList().Front(); element != nil; element = element.Next() {
		sentences[i] = element.Value.(string)
		i++
	}
	return &CStorage{id, conv.GetPrompt(), sentences}
}
func PlainText(conv Conversation) string {
	isGpt3 := false
	switch conv.(type) {
	case *GPT3Conversation:
	default:
		isGpt3 = true
	}
	builder := new(strings.Builder)
	builder.WriteString(conv.GetPrompt())
	builder.WriteString("\n")
	index := 0
	for element := conv.GetSentenceList().Front(); element != nil; element = element.Next() {
		if index%2 == 0 {
			if isGpt3 {
				builder.WriteString("\n")
				builder.WriteString(conv.GetHumanName())
				builder.WriteString(": ")
			} else {
				builder.WriteString(conv.GetHumanName())
			}
		} else {
			if isGpt3 {
				builder.WriteString("\n")
				builder.WriteString(conv.GetAIName())
				builder.WriteString(": ")
			} else {
				builder.WriteString(conv.GetAIName())
			}
		}
		builder.WriteString(element.Value.(string))
		index++
	}
	return builder.String()
}

func Rollback(conv *Conversation) error {
	sentenceList := (*conv).GetSentenceList()
	length := sentenceList.Len()
	if length < 2 {
		return gpterror.ChatgptError{Msg: "Failed to rollback, dialog is empty"}
	} else if length%2 != 0 {
		return gpterror.ChatgptError{Msg: "Failed to rollback, cant rollback while AI is thinking"}
	} else {
		sentenceList.Remove(sentenceList.Back())
		sentenceList.Remove(sentenceList.Back())
		return nil
	}
}

func CreateQuickConversation(prompt string, sentenceList *list.List, stream bool) *GPT3Conversation {
	return &GPT3Conversation{Prompt: prompt, SentenceList: sentenceList, AIAnswered: true, LastModify: time.Now().Unix(), RequestSettings: QuickChatSettings, AIName: "\nAI: ", HumanName: "\nHuman: "}
}

func CreateQuickConversationTurbo(prompt string, sentenceList *list.List, stream bool) *TurboConversation {
	return &TurboConversation{Prompt: prompt, SentenceList: sentenceList, AIAnswered: true, LastModify: time.Now().Unix(), AIName: "assistant", HumanName: "user"}
}

func CreateQuickConversationGPT4(prompt string, sentenceList *list.List, stream bool) *GPT4Conversation {
	return &GPT4Conversation{Prompt: prompt, SentenceList: sentenceList, AIAnswered: true, LastModify: time.Now().Unix(), AIName: "assistant", HumanName: "user"}
}
