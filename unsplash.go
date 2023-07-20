package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UnsplashRandom struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Color       string `json:"color"`
	Description string `json:"description"`
	URLs        struct {
		Raw     string `json:"raw"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Full    string `json:"full"`
		Thumb   string `json:"thumb"`
		S3      string `json:"small-s3"`
	}
	Links struct {
		Self             string `json:"self"`
		HTML             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	User struct {
		ID        string `json:"id"`
		UpdatedAt string `json:"updated_at"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		Links     struct {
			Self      string `json:"self"`
			HTML      string `json:"html"`
			Photos    string `json:"photos"`
			Likes     string `json:"likes"`
			Portfolio string `json:"portfolio"`
		} `json:"links"`
	} `json:"user"`
}

type UnsplashDownloadLoc struct {
	URL string `json:"url"`
}

func UnsplashImageFromApi(query string) *UnsplashRandom {
	resp, err := http.Get(fmt.Sprintf("https://api.unsplash.com/photos/random?query=%v&client_id=%v", query, *UnsplashToken))
	if err != nil {
		ChannelLog(fmt.Sprintf("An Error Occured %v", err))
	}

	var unsplash *UnsplashRandom

	err = json.NewDecoder(resp.Body).Decode(&unsplash)
	if err != nil {
		var invalid *UnsplashRandom
		return invalid
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return unsplash
}
