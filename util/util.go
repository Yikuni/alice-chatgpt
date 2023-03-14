package util

import (
	"alice-chatgpt/global"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func PostHeader(_url string, msg []byte, headers map[string]string) (string, error) {
	var client *http.Client
	if global.Proxy != "" {
		proxyUrl, err := url.Parse(global.Proxy)
		if err == nil {
			client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
		} else {
			fmt.Println("Proxy Failed, using non-proxy")
		}
	}
	if client == nil {
		client = &http.Client{}
	}
	req, err := http.NewRequest("POST", _url, strings.NewReader(string(msg)))
	if err != nil {
		return "", err
	}
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(response.Body)
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

func HexBuffToString(buff []byte) string {
	var ret string
	for _, value := range buff {
		str := strconv.FormatUint(uint64(value), 16)
		if len([]rune(str)) == 1 {
			ret = ret + "0" + str
		} else {
			ret = ret + str
		}
	}
	return ret
}
