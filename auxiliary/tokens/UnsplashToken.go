package tokens

import (
	"log"
	"os"
)

var UnsplashToken string

func GetUnsplashToken() {
	UnsplashToken = ReadTokenFromFile("unsplash-token.txt")

	if UnsplashToken != "" {
		log.Println("Token read from file")
	} else {
		log.Println("Token not read from file, fetching from env")
		UnsplashToken = os.Getenv("UNSPLASH_TOKEN")
	}
}
