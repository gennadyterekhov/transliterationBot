package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"transliteration_bot/controllers"
	"transliteration_bot/telegramAPI"
	"transliteration_bot/webhook"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/webhook", webhook.Webhook)
	// fmt.Println(os.Getenv("webhook_url"))

	fmt.Println("[current webhook info]:[\n", telegramAPI.GetWebhookInfo(), "\n]")

	telegramAPI.SetWebhook(os.Getenv("webhook_url"))

	if !telegramAPI.IsWebhookSet() {
		telegramAPI.SetWebhook(os.Getenv("webhook_url"))
	}

	fmt.Printf("server started on %s\n", os.Getenv("port"))
	http.ListenAndServe(":"+os.Getenv("port"), nil)
}
