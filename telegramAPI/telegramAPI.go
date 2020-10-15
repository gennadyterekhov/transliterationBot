package telegramAPI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"transliteration_bot/helper"
)

type getWebhookInfoResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		URL                  string `json:"url"`
		HasCustomCertificate bool   `json:"has_custom_certificate"`
		PendingUpdateCount   int    `json:"pending_update_count"`
	} `json:"result"`
}

func telegramApiRequest(methodName string, params string) string {
	// local
	// var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/", config.GetConfig().APIKey)
	// prod
	var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("api_key"))

	var response *http.Response
	var err error
	if params == "" {
		response, err = http.Get(fullUrl + methodName)
	} else {
		response, err = http.Get(fullUrl + methodName + "?" + url.PathEscape(params))
	}
	helper.Check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return string(body)
}

func GetMe(w http.ResponseWriter, req *http.Request) {
	var response string = telegramApiRequest("getMe", "")
	fmt.Fprintf(w, response)
}

func SetWebhook(webhookUrl string) string {
	// local
	// var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", config.GetConfig().APIKey)
	// prod
	var fullUrl string = fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook", os.Getenv("api_key"))
	response, err := http.PostForm(
		fullUrl,
		url.Values{"url": {webhookUrl}},
	)

	helper.Check(err)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	return string(body)
}

func SendMessage(chatId int, text string) string {
	var params string = fmt.Sprintf("chat_id=%d&text=%s", chatId, text)
	var response string = telegramApiRequest("sendMessage", params)
	return response
}

// сюда приходят хуки, эта функция запускает обработчик сообщений

func GetWebhookInfo() string {
	var response string = telegramApiRequest("getWebhookInfo", "")
	return response
}

func IsWebhookSet() bool {
	var response string = GetWebhookInfo()
	var webhookInfo getWebhookInfoResponse = getWebhookInfoResponse{}
	if err := json.Unmarshal([]byte(response), &webhookInfo); err != nil {
		panic(err)
	}
	if webhookInfo.Result.URL == "" {
		return false
	}
	return true
}
