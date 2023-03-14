package conversation

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"time"
)

type TurboConversation struct {
	Prompt       string     // prompt
	SentenceList *list.List // 对话表, 偶数是人类
	AIAnswered   bool       // AI是否完成回复
	LastModify   int64      // 上次回复的时间戳
	AIName       string
	HumanName    string
	Conversation
}

type RoleContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TurboRequest struct {
	Model    string        `json:"model"`
	Messages []RoleContent `json:"messages"`
}

func (conv *TurboConversation) GetHumanName() string {
	return conv.HumanName
}
func (conv *TurboConversation) GetAIName() string {
	return conv.AIName
}
func (conv *TurboConversation) GetLastModify() int64 {
	return conv.LastModify
}
func (conv *TurboConversation) SetLastModify(lastModify int64) {
	conv.LastModify = lastModify
}
func (conv *TurboConversation) GetAIAnswered() bool {
	return conv.AIAnswered
}
func (conv *TurboConversation) SetAIAnswered(AIAnswered bool) {
	conv.AIAnswered = AIAnswered
}
func (conv *TurboConversation) GetSentenceList() *list.List {
	return conv.SentenceList
}
func (conv *TurboConversation) GetPrompt() string {
	return conv.Prompt
}
func (conv *TurboConversation) RequestBody() ([]byte, error) {
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
	jsonString, err := json.Marshal(&TurboRequest{Model: "gpt-3.5-turbo-0301", Messages: msgArray})
	fmt.Println(string(jsonString))
	if err != nil {
		return nil, err
	} else {
		return jsonString, nil
	}
}
func (conv *TurboConversation) SolveResponse(jsonObject *gabs.Container) string {
	fmt.Println(jsonObject.String())
	return jsonObject.S("choices", "0", "message", "content").Data().(string)
}

func CreateTuborConversation(prompt string, AIName string, HumanName string) *TurboConversation {
	return &TurboConversation{Prompt: prompt, SentenceList: list.New(), AIAnswered: true, LastModify: time.Now().Unix(), AIName: AIName, HumanName: HumanName}
}
