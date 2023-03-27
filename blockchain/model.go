package blockchain

import (
	"time"
)

type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Transactions []*Transaction
	Timestamp    time.Time
	Pow          int
}
type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}
type Transaction struct {
	From   Block
	To     []Block
	Amount int
}
