package util

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func PostHeader(url string, msg []byte, headers map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(msg)))
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
