package lib

import (
	"log"
	"os"
)

func Debug(message string) {
	if os.Getenv("env") == "dev" {
		log.Println(message)

	}

}
