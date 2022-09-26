package lib

import "log"

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)

	}
}

func ErrorByte(b []byte, e error) []byte {
	CheckError(e)
	return b
}
