package conversation

import (
	"container/list"
	"encoding/json"
	"github.com/Jeffail/gabs/v2"
	"time"
)

type GPT4Conversation struct {
	Prompt       string     // prompt
	SentenceList *list.List // 对话表, 偶数是人类
	AIAnswered   bool       // AI是否完成回复
	LastModify   int64      // 上次回复的时间戳
	AIName       string
	HumanName    string
	Conversation
}

func (conv *GPT4Conversation) GetHumanName() string {
	return conv.HumanName
}
func (conv *GPT4Conversation) GetAIName() string {
	return conv.AIName
}
func (conv *GPT4Conversation) GetLastModify() int64 {
	return conv.LastModify
}
func (conv *GPT4Conversation) SetLastModify(lastModify int64) {
	conv.LastModify = lastModify
}
func (conv *GPT4Conversation) GetAIAnswered() bool {
	return conv.AIAnswered
}
func (conv *GPT4Conversation) SetAIAnswered(AIAnswered bool) {
	conv.AIAnswered = AIAnswered
}
func (conv *GPT4Conversation) GetSentenceList() *list.List {
	return conv.SentenceList
}
func (conv *GPT4Conversation) GetPrompt() string {
	return conv.Prompt
}
func (conv *GPT4Conversation) RequestBody() ([]byte, error) {
	arrayLength := conv.SentenceList.Len() + 1
	msgArray := make([]RoleContent, arrayLength)
	msgArray[0] = RoleContent{Role: "system", Content: conv.Prompt}
	index := 1
	for element := conv.SentenceList.Front(); element != nil; element = element.Next() {
		if index%2 == 0 {
			msgArray[index] = RoleContent{Role: conv.HumanName, Content: element.Value.(string)}
		} else {
			msgArray[index] = RoleContent{Role: conv.AIName, Content: element.Value.(string)}
		}
		index++
	}
	jsonString, err := json.Marshal(&TurboRequest{Model: "gpt-4", Messages: msgArray})
	if err != nil {
		return nil, err
	} else {
		return jsonString, nil
	}
}
func (conv *GPT4Conversation) SolveResponse(jsonObject *gabs.Container) string {
	return jsonObject.S("choices", "0", "message", "content").Data().(string)
}

func CreateGPT4Conversation(prompt string, AIName string, HumanName string) *GPT4Conversation {
	return &GPT4Conversation{Prompt: prompt, SentenceList: list.New(), AIAnswered: true, LastModify: time.Now().Unix(), AIName: AIName, HumanName: HumanName}
}
