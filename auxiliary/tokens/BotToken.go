package tokens

import (
	"log"
	"os"
)

var BotToken string

func GetBotToken() {
	BotToken = ReadTokenFromFile("token-dev.txt")

	if BotToken != "" {
		log.Println("Token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		BotToken = os.Getenv("TOKEN")
	}
}
