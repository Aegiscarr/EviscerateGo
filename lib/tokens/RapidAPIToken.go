package tokens

import (
	"log"
	"os"
)

var RapidSzToken string

func GetRapidApiToken() {
	RapidSzToken = ReadTokenFromFile("cfg/tokens/rapid-sz-token.txt")

	if RapidSzToken != "" {
		log.Println("RapidAPI Token read from file")
	} else {
		log.Println("RapidAPI Token not read from file, fetching from env")
		RapidSzToken = os.Getenv("RAPIDAPI_TOKEN")
	}
}
