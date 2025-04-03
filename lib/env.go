package lib

import (
	"os"
	"slices"
)

func IsLocal() bool {
	return slices.Contains([]string{"local", ""}, os.Getenv("env"))
}
