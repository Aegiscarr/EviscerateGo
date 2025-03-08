package conf

import "log"

var AdminId string

func GetAdminId() {
	AdminId = ReadTokenFromFile("cfg/admin-id.txt")

	if AdminId != "" {
		log.Println("AdminId read from file")
	} else {
		log.Println("No AdminId specified, admin-only commands will be unusable!")
	}
}
