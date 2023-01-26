// @author zhuweitung 2023/1/26
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   string `json:"error_code"`
	Description string `json:"description"`
	Result      struct {
		MessageId int `json:"message_id"`
		From      struct {
			Id        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

func Send(message string) {
	notification := CONFIG.Notification
	if notification == (Notification{}) {
		log.Println("未检测到消息通知配置，跳过消息通知")
		return
	}
	telegram := notification.Telegram
	if telegram != (Telegram{}) &&
		telegram.BotToken != "" &&
		telegram.ChatId != "" {
		// 发送电报消息
		apiUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegram.BotToken)
		data := make(url.Values)
		data["chat_id"] = []string{telegram.ChatId}
		data["text"] = []string{message}
		sendTelegram(getHttpClient(telegram.HttpProxy), apiUrl, data)
	}
}

// 获取http客户端
func getHttpClient(httpProxy string) (httpclient http.Client) {
	httpclient = http.Client{
		// 设置超时
		Timeout: time.Duration(time.Second * 30),
	}
	if httpProxy != "" {
		// 设置http代理
		proxyURL, _ := url.Parse(httpProxy)
		httpclient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}
	return httpclient
}

func sendTelegram(httpclient http.Client, url string, data url.Values) {
	resp, err := httpclient.PostForm(url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var tgResp TelegramResponse
	if err := json.Unmarshal(body, &tgResp); err != nil {
		panic(err)
	}
	if tgResp.Ok {
		log.Println("推送电报消息成功")
	} else {
		log.Fatalln("推送电报消息失败: ", tgResp.Description)
	}
}
