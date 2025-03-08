package tokens

import (
	"log"
	"os"
)

var UploaderToken string

func GetUploaderToken() {
	UploaderToken = ReadTokenFromFile("cfg/tokens/uploader-token.txt")

	if UploaderToken != "" {
		log.Println("Uploader Token read from file")
	} else {
		log.Println("Uploader Token not read from file, fetching from env")
		UploaderToken = os.Getenv("UPLOADER_TOKEN")
	}
}
