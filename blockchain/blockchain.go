package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	lib "github.com/wopta/goworkspace/lib"
	models "github.com/wopta/goworkspace/models"
)

func Blockchain(w http.ResponseWriter, r *http.Request) {
	lib.EnableCors(&w, r)

}

func calculateHashf(w http.ResponseWriter, r *http.Request) {}
func calculateHash(b models.Block) string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.Pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func mine(difficulty int, b models.Block) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Pow++
		b.Hash = calculateHash(b)
	}
}

func addBlock(from, to string, amount float64, b models.Blockchain) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := models.Block{
		Data:         blockData,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	mine(b.Difficulty, newBlock)
	b.Chain = append(b.Chain, newBlock)
}
func isValid(b models.Blockchain) bool {
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		if currentBlock.Hash != calculateHash(currentBlock) || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}
func CreateBlockchain(difficulty int) models.Blockchain {
	genesisBlock := models.Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}
	return models.Blockchain{
		GenesisBlock: genesisBlock,
		Chain:        []models.Block{genesisBlock},
		Difficulty:   difficulty,
	}
}
