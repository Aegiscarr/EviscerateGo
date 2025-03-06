package api

import (
	"EviscerateGo/auxiliary/structs"
	"EviscerateGo/auxiliary/tokens"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetRapidAPICall(parsedQ string, typeQ string, RapidSzToken string) *structs.RapidSzResponse {

	var (
		szresponse  *structs.RapidSzResponse
		parsedQuery = parsedQ
		queryType   = typeQ
		url         string
	)

	url = "https://spotify23.p.rapidapi.com/search/?q=" + parsedQuery + "&type=" + queryType + "&offset=0&limit=1&numberOfTopResults=1"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", tokens.RapidSzToken)
	req.Header.Add("X-RapidAPI-Host", "spotify23.p.rapidapi.com")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error while sending request: ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	fmt.Println(json.Valid(body))
	fmt.Println(json.Unmarshal(body, &szresponse))

	return szresponse
}
