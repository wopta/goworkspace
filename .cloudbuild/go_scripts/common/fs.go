package common

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetGoMod(dir, module string) []byte {
	path, err := filepath.Abs(fmt.Sprintf("%s/%s/go.mod", dir, module))
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}
