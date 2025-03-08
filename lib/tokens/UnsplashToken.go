package tokens

import (
	"log"
	"os"
)

var UnsplashToken string

func GetUnsplashToken() {
	UnsplashToken = ReadTokenFromFile("cfg/tokens/unsplash-token.txt")

	if UnsplashToken != "" {
		log.Println("Unsplash Token read from file")
	} else {
		log.Println("Unsplash Token not read from file, fetching from env")
		UnsplashToken = os.Getenv("UNSPLASH_TOKEN")
	}
}
