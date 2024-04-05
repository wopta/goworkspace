package lib

import "os"

func IsLocal() bool {
	return os.Getenv("env") == "local"
}
