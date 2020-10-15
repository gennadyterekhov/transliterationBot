package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"transliteration_bot/helper"
	"transliteration_bot/telegramAPI"
	"transliteration_bot/transliterate"
)

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

func getMessageObj(s string) messageWebhook {
	var messageObj messageWebhook = messageWebhook{}
	if err := json.Unmarshal([]byte(s), &messageObj); err != nil {
		panic(err)
	}
	return messageObj
}
func getResponseByMessageText(messageText string) string {
	if messageText == "/start" {
		return "Hello!\nI'm transliterator bot.\nI can transliterate russian cyrillic into latin."
	}
	return transliterate.Transliterate(messageText)
}
func Webhook(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "this is where bot lives. move along")

	body, err := ioutil.ReadAll(req.Body)
	helper.Check(err)
	fmt.Printf("\n[received webhook] %s\n", string(body))

	var message messageWebhook = messageWebhook{}
	message = getMessageObj(string(body))

	fmt.Printf("\n[got message, generating answer]: message.Message.Text is %s\n", message.Message.Text)
	var response string = getResponseByMessageText(message.Message.Text)

	fmt.Printf("\n[sending message]: response var is %s\n", response)
	telegramAPI.SendMessage(message.Message.From.ID, response)

}
