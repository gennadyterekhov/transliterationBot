package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type config struct {
	Port       int    `json:"port"`
	APIKey     string `json:"api_key"`
	WebhookUrl string `json:"webhook_url"`
	Cyrillic   string `json:"cyrillic"`
	Latin      string `json:"latin"`
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

func getWebhookUrlFromConfig() string {
	return getConfig().WebhookUrl
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

func telegramApiRequest(methodName string, params string) string {
	var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/", getApiKeyFromConfig())
	var response *http.Response
	var err error
	if params == "" {
		response, err = http.Get(fullUrl + methodName)
	} else {
		response, err = http.Get(fullUrl + methodName + "?" + url.PathEscape(params))
	}
	check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	// fmt.Printf("return from telegramApiRequest: %s\n", string(body))
	return string(body)
}

// func telegramApiRequestPost(methodName string) string {
// 	var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/", getApiKeyFromConfig())

// 	response, err := http.PostForm(
// 		fullUrl+url.PathEscape(methodName),
// 		url.Values{"url": {"Value"}},
// 	)

// 	check(err)
// 	defer response.Body.Close()

// 	body, err := ioutil.ReadAll(response.Body)

// 	return string(body)
// }

func getMe(w http.ResponseWriter, req *http.Request) {
	var response string = telegramApiRequest("getMe", "")
	fmt.Fprintf(w, response)
}

func setWebhook(webhookUrl string) string {
	var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", getApiKeyFromConfig())
	response, err := http.PostForm(
		fullUrl,
		url.Values{"url": {webhookUrl}},
	)

	check(err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return string(body)
}

func sendMessage(chatId int, text string) string {
	var params string = fmt.Sprintf("chat_id=%d&text=%s", chatId, text)
	var response string = telegramApiRequest("sendMessage", params)
	return response
}

func transliterate(msg string) string {
	var cyrillic string = getConfig().Cyrillic
	var latin string = getConfig().Latin

	const numberOfChars int = 66

	var cyrillicRunes [numberOfChars]rune
	var latinRunes [numberOfChars]rune
	var trans string = ""

	var correspondence map[rune]rune = make(map[rune]rune)

	var j uint8 = 0
	for _, runeValue := range cyrillic {
		cyrillicRunes[j] = runeValue
		j += 1
	}
	j = 0
	for _, runeValue := range latin {
		latinRunes[j] = runeValue
		j += 1
	}

	for i := 0; i < numberOfChars; i += 1 {
		correspondence[cyrillicRunes[i]] = latinRunes[i]
	}

	for _, runeValue := range msg {
		if latinChar, ok := correspondence[runeValue]; ok {
			//do something here
			trans += string(latinChar)
		} else {
			trans += string(runeValue)
		}

	}

	return trans

}

func getResponseByMessageText(messageText string) string {
	if messageText == "/start" {
		return "Hello! I'm transliterator bot."
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

	fmt.Printf("\n[got message, generating answer]: message.Message.Text is %s\n", message.Message.Text)
	var response string = getResponseByMessageText(message.Message.Text)

	fmt.Printf("\n[sending message]: response var is %s\n", response)
	sendMessage(message.Message.From.ID, response)

}

func getWebhookInfo() string {
	var response string = telegramApiRequest("getWebhookInfo", "")
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

func router() {
	http.HandleFunc("/", index)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/webhook", webhook)
	http.HandleFunc("/getMe", getMe)
}

func main() {
	router()

	fmt.Println("[current webhook info]:[\n", getWebhookInfo(), "\n]")

	setWebhook(getWebhookUrlFromConfig())
	if !isWebhookSet() {
		setWebhook(getWebhookUrlFromConfig())
	}

	fmt.Printf("server started on %d\n", getPortFromConfig())

	http.ListenAndServe(":"+fmt.Sprint(getPortFromConfig()), nil)
}
