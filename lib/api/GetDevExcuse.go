package api

import (
	"EviscerateGo/lib/structs"
	"fmt"

	"encoding/json"
	"io"
	"net/http"
)

func GetDevExcuse() *structs.RandomDevExcuse {
	resp, err := http.Get(("https://api.tabliss.io/v1/developer-excuses"))
	if err != nil {
		fmt.Printf("An error occured: %v", err)
	}
	var randomexcuse *structs.RandomDevExcuse
	err = json.NewDecoder(resp.Body).Decode(&randomexcuse)
	if err != nil {
		var invalid *structs.RandomDevExcuse
		return invalid
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	return randomexcuse
}
