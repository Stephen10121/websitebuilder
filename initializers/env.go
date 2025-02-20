package initializers

import (
	"log"
	"myapp/funcs"
	"os"
)

var AdminUsername string
var AdminPassword string
var AdminSalt string

func SetupEnv() {
	adminUsername := os.Getenv("ADMIN_USERNAME")

	if len(adminUsername) != 0 {
		AdminUsername = adminUsername
	} else {
		log.Fatalln("ADMIN_USERNAME env variable is not set.")
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if len(adminPassword) != 0 {
		AdminPassword = adminPassword
	} else {
		log.Fatalln("ADMIN_PASSWORD env variable is not set.")
	}

	AdminSalt = funcs.RandStringBytes(8)
}
