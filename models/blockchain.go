package models

import (
	"time"
)

type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Pow          int
}
type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}
