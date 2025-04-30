package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Block struct {
	Index        uint64 `json:"index"`
	Timestamp    int64  `json:"timestamp"`
	ShipmentData string `json:"shipment_data"`
	PrevHash     string `json:"prev_hash"`
	Validator    string `json:"validator"`
	Hash         string `json:"hash"`
}

func (b *Block) ComputeHash() string {
	record :=
		fmt.Sprintf("%d%d%s%s%s",
			b.Index,
			b.Timestamp,
			b.ShipmentData, b.PrevHash, b.Validator)
	sum := sha256.Sum256([]byte(record))
	return hex.EncodeToString(sum[:])
}

var AuthorizedValidators = []string{"validator1_pubkey", "validator2_pubkey"}

func isAuthorized(val string) bool {
	for _, v := range AuthorizedValidators {
		if v == val {
			return true
		}
	}
	return false
}

type Blockchain struct {
	Chain []Block
}

func NewBlockchain() *Blockchain {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		ShipmentData: "genesis",
		PrevHash:     "0",
		Validator:    "system",
	}
	genesis.Hash = genesis.ComputeHash()
	return &Blockchain{Chain: []Block{genesis}}
}

func (bc *Blockchain) AddBlock(data, validator string) error {
	if !isAuthorized(validator) {
		return errors.New("validator not authorized")
	}
	prev := bc.Chain[len(bc.Chain)-1]
	blk := Block{
		Index:        prev.Index + 1,
		Timestamp:    time.Now().Unix(),
		ShipmentData: data,
		PrevHash:     prev.Hash,
		Validator:    validator,
	}
	blk.Hash = blk.ComputeHash()
	bc.Chain = append(bc.Chain, blk)
	return nil
}

func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		curr, prev := bc.Chain[i], bc.Chain[i-1]
		if curr.PrevHash != prev.Hash || curr.Hash != curr.ComputeHash() {
			return false
		}
		if !isAuthorized(curr.Validator) {
			return false
		}
	}
	return true
}

func main() {
	bc := NewBlockchain()
	fmt.Println("Genesis:", bc.Chain[0])

	// Успешная попытка.
	err := bc.AddBlock(`{"id":"S1","status":"shipped"}`,
		"validator1_pubkey")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Неуспешная попытка.
	err = bc.AddBlock(`{"id":"S2","status":"pending"}`, "bad_validator")
	fmt.Println("AddBlock with bad_validator:", err)

	// Вывод всей цепочки.
	out, _ := json.MarshalIndent(bc.Chain, "", " ")
	fmt.Println(string(out))

	// Проверка валидности.
	fmt.Println("Chain valid?", bc.IsValid())
}
