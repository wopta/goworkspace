package lib

import (
	"log"
	"os"
)

func Debug(v ...interface{}) {
	if os.Getenv("env") == "dev" {
		log.Println(v)

	}

}
