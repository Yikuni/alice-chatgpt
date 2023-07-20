package conversation

import (
	"alice-chatgpt/global"
	"alice-chatgpt/gpterror"
	"alice-chatgpt/util"
	"bufio"
	"bytes"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type APIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func CallStreamAPI(conv Conversation, question string, c *gin.Context) (string, error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Allow-Origin", "*")
	if question == "" {
		return "", gpterror.ChatgptError{Msg: "Question can not be empty"}
	}
	if !conv.GetAIAnswered() {
		return "", gpterror.ChatgptError{Msg: "AI is thinking"}
	}
	conv.SetAIAnswered(false)
	conv.GetSentenceList().PushBack(question)
	key := Key()
	var _url string
	switch conv.(type) {
	case *GPT3Conversation:
		_url = "https://api.openai.com/v1/completions"
	case *TurboConversation:
		_url = "https://api.openai.com/v1/chat/completions"
	case *GPT4Conversation:
		_url = "https://api.openai.com/v1/chat/completions"
	}
	jsonString, err := conv.RequestBody()
	if err != nil {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		fmt.Println(err)
		return "", err
	}
	req, err := http.NewRequest("POST", _url, bytes.NewBuffer(jsonString))
	if err != nil {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", key))

	//result, err := util.PostHeader(_url, jsonString, headers)
	client := util.GetClient()
	resp, err := client.Do(req)
	defer func() {
		conv.SetAIAnswered(true)
		conv.SetLastModify(time.Now().Unix())
		err := resp.Body.Close()
		if err != nil {
			fmt.Println("Failed to close response")
		}
	}()
	if err != nil {
		conv.GetSentenceList().Remove(conv.GetSentenceList().Back())
		if err.Error() == "Post \"https://api.openai.com/v1/completions\": net/http: invalid header field value for \"Authorization\"" && global.AutoRemoveErrorKeys {
			findAndRemoveKey(key)
		}
		fmt.Println(err)
		return "", err
	}
	reader := bufio.NewReader(resp.Body)
	var answer string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取时出错")
			return "", err
		}
		if strings.HasPrefix(line, "data:") {
			line = strings.TrimPrefix(line, "data:")
			//fmt.Println("Received data:", line)

			if strings.TrimSpace(line) == "[DONE]" {
				conv.GetSentenceList().PushBack(answer)
				msg := fmt.Sprintf("data: %s\n\n", "[DONE]")
				_, err := c.Writer.Write([]byte(msg))
				if err != nil {
					return "", err
				}
				c.Writer.Flush()
				return answer, nil
			}

			jsonObject, err := gabs.ParseJSON([]byte(line))
			if err != nil {
				fmt.Println("解析时出错, line为" + line)
				return "", err
			}
			if jsonObject.Exists("choices") {
				choices := jsonObject.S("choices").Children()
				if len(choices) > 0 {
					data := choices[0].S("delta", "content").Data()
					if data != nil {
						content := data.(string)
						answer += content
						msg := fmt.Sprintf("data: %s\n\n", content)
						_, err := c.Writer.Write([]byte(msg))
						if err != nil {
							return "", err
						}
						c.Writer.Flush()
					}
				}
			}
		}
	}
}
