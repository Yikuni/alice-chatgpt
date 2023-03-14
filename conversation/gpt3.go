package conversation

import (
	"container/list"
	"encoding/json"
	"github.com/Jeffail/gabs/v2"
	"time"
)

type GPT3Conversation struct {
	Prompt          string     // prompt
	SentenceList    *list.List // 对话表, 偶数是人类
	AIAnswered      bool       // AI是否完成回复
	LastModify      int64      // 上次回复的时间戳
	RequestSettings *RequestSettings
	AIName          string
	HumanName       string
	Conversation
}

func (conv *GPT3Conversation) GetHumanName() string {
	return conv.HumanName
}
func (conv *GPT3Conversation) GetAIName() string {
	return conv.AIName
}
func (conv *GPT3Conversation) GetLastModify() int64 {
	return conv.LastModify
}
func (conv *GPT3Conversation) SetLastModify(lastModify int64) {
	conv.LastModify = lastModify
}
func (conv *GPT3Conversation) GetAIAnswered() bool {
	return conv.AIAnswered
}
func (conv *GPT3Conversation) SetAIAnswered(AIAnswered bool) {
	conv.AIAnswered = AIAnswered
}
func (conv *GPT3Conversation) GetSentenceList() *list.List {
	return conv.SentenceList
}
func (conv *GPT3Conversation) GetPrompt() string {
	return conv.Prompt
}
func CreateCustomConversation(prompt string, settings *RequestSettings, AIName string, HumanName string) *GPT3Conversation {
	//return &GPT3Conversation{prompt, list.New(), true, time.Now().Unix(), settings, AIName, HumanName}
	return &GPT3Conversation{Prompt: prompt, SentenceList: list.New(), AIAnswered: true, LastModify: time.Now().Unix(), RequestSettings: settings, AIName: AIName, HumanName: HumanName}
}

func CreateQuickConversation(prompt string, sentenceList *list.List) *GPT3Conversation {
	return &GPT3Conversation{Prompt: prompt, SentenceList: sentenceList, AIAnswered: true, LastModify: time.Now().Unix(), RequestSettings: QuickChatSettings, AIName: "\nAI: ", HumanName: "\nHuman: "}
}

// RequestBody
/**
Get Request body
Question should be pushed into sentence list before call this function
*/
func (conv *GPT3Conversation) RequestBody() ([]byte, error) {
	prompt := PlainText(conv) + conv.AIName
	request := ChatgptRequest{
		Prompt:          prompt,
		RequestSettings: *conv.RequestSettings,
	}
	jsonString, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	} else {
		return jsonString, nil
	}
}

func (conv *GPT3Conversation) SolveResponse(jsonObject *gabs.Container) string {
	return jsonObject.S("choices", "0", "text").Data().(string)
}
