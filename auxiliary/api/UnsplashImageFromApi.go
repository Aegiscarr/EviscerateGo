package api

import (
	"EviscerateGo/auxiliary/structs"
	"EviscerateGo/auxiliary/tokens"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnsplashImageFromApi(query string) *structs.UnsplashRandom {
	resp, err := http.Get(fmt.Sprintf("https://api.unsplash.com/photos/random?query=%v&client_id=%v", query, tokens.UnsplashToken))
	if err != nil {
		fmt.Printf("An Error Occured %v", err)
	}

	var unsplash *structs.UnsplashRandom

	err = json.NewDecoder(resp.Body).Decode(&unsplash)
	if err != nil {
		var invalid *structs.UnsplashRandom
		return invalid
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return unsplash
}
