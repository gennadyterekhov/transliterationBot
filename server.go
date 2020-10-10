package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type config struct {
	Port     int    `json:"port"`
	APIKey   string `json:"api_key"`
	Cyrillic string `json:"cyrillic"`
	Latin    string `json:"latin"`
}

type getWebhookInfoResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		URL                  string `json:"url"`
		HasCustomCertificate bool   `json:"has_custom_certificate"`
		PendingUpdateCount   int    `json:"pending_update_count"`
	} `json:"result"`
}

type messageWebhook struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMessageObj(s string) messageWebhook {
	var messageObj messageWebhook = messageWebhook{}
	if err := json.Unmarshal([]byte(s), &messageObj); err != nil {
		panic(err)
	}
	return messageObj
}

func getConfig() config {
	configStr, err := ioutil.ReadFile("config.json")
	check(err)

	var configObj config = config{}
	if err := json.Unmarshal([]byte(configStr), &configObj); err != nil {
		panic(err)
	}
	return configObj
}

func getApiKeyFromConfig() string {
	return getConfig().APIKey
}

func getPortFromConfig() int {
	return getConfig().Port
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "you are at the index\n")
}

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func telegramApiRequest(methodName string) string {
	var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/", getApiKeyFromConfig())
	response, err := http.Get(fullUrl + url.PathEscape(methodName))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	// fmt.Printf("return from telegramApiRequest: %s\n", string(body))
	return string(body)
}

func getMe(w http.ResponseWriter, req *http.Request) {
	var response string = telegramApiRequest("getMe")
	fmt.Fprintf(w, response)
}

func setWebhook(fullUrl string) string {
	var urlWithParams string = fmt.Sprintf("setWebhook?url=%s", fullUrl)
	// тут надо сделать чтобы пост отправлял
	var response string = telegramApiRequest(urlWithParams)
	return response
}

func sendMessage(chatId int, text string) string {
	var urlWithParams string = fmt.Sprintf("sendMessage?chat_id=%d&text=%s", chatId, text)
	var response string = telegramApiRequest(urlWithParams)
	return response
}

func transliterate(msg string) string {
	var cyrillic string = getConfig().Cyrillic
	var latin string = getConfig().Latin
	var trans string = ""

	for i := 0; i < len(msg); i += 1 {
		// if msg[i] in cyrillic
		// trans += latin[cyrillic.indexOf(msg[i])]
		// else trans += msg[i]
		var charIndex int = strings.Index(cyrillic, string(msg[i]))
		if charIndex != -1 {
			trans += string(latin[charIndex])
		} else {
			trans += string(msg[i])
		}
	}

	return trans
}

func getResponseByMessageText(messageText string) string {
	if messageText == "/start" {
		return "Hello!\nI'm transliterator bot."
	}
	return transliterate(messageText)
}

// сюда приходят хуки, эта функция запускает обработчик сообщений
func webhook(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "this is where bot lives. move along")

	body, err := ioutil.ReadAll(req.Body)
	check(err)
	fmt.Printf("\n[received webhook] %s\n", string(body))

	var message messageWebhook = messageWebhook{}
	message = getMessageObj(string(body))
	var response string = getResponseByMessageText(message.Message.Text)

	sendMessage(message.Message.From.ID, response)

}

// func getUpdates() string {
// 	var response string = telegramApiRequest("getUpdates")
// 	return response
// }

func getWebhookInfo() string {
	var response string = telegramApiRequest("getWebhookInfo")
	return response
}

func isWebhookSet() bool {
	var response string = getWebhookInfo()
	var webhookInfo getWebhookInfoResponse = getWebhookInfoResponse{}
	if err := json.Unmarshal([]byte(response), &webhookInfo); err != nil {
		panic(err)
	}
	if webhookInfo.Result.URL == "" {
		return false
	}
	return true
}

func main() {
	// examples
	http.HandleFunc("/", index)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	//
	//
	http.HandleFunc("/webhook", webhook)
	http.HandleFunc("/getMe", getMe)

	fmt.Println(getWebhookInfo())

	setWebhook("https://85d449ffee77.ngrok.io/webhook")
	fmt.Println(getWebhookInfo())
	if !isWebhookSet() {
		// setWebhook("")
		setWebhook("https://85d449ffee77.ngrok.io/webhook")
	}

	// fmt.Println("[webhook info] ", getWebhookInfo())
	// sendMessage("@gennadyterekhov", "test message")
	fmt.Printf("server started on %d\n", getPortFromConfig())

	http.ListenAndServe(":"+fmt.Sprint(getPortFromConfig()), nil)
}
