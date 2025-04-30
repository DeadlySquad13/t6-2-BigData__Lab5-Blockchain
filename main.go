package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Shipment struct {
	Id     string
	From   string
	To     string
	Status string
}

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

func (bc *Blockchain) AddShipment(shipment Shipment, validator string) error {
	shipmentJson, err := json.Marshal(shipment)

	if err != nil {
		return err
	}

	return bc.AddBlock(string(shipmentJson), validator)
}

func (bc *Blockchain) AddShipments(shipments []Shipment, validator string) error {
	for _, shipment := range shipments {
		err := bc.AddShipment(shipment, validator)
		if err != nil {
			return err
		}
	}

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
	shipments_company1 := []Shipment{
		{"S1", "Warehouse A", "Client X", "shipped"},
		{"S2", "Warehouse B", "Client Y", "in_transit"},
		{"S3", "Warehouse C", "Client Z", "delivered"},
		{"S4", "Warehouse A", "Client W", "in_transit"},
	}

	shipments_company2 := []Shipment{
		{"S5", "Warehouse D", "Client W", "in_transit"},
		{"S6", "Warehouse E", "Client X", "delivered"},
	}

	shipments_company3 := []Shipment{
		{"S7", "Warehouse A", "Client X", "shipped"},
		{"S8", "Warehouse A", "Client Y", "in_transit"},
		{"S9", "Warehouse A", "Client Z", "in_transit"},
	}

	bc := NewBlockchain()
	fmt.Println("Genesis:", bc.Chain[0])

	// Успешная попытка.
	err := bc.AddShipments(shipments_company1,
		AuthorizedValidators[0])
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Успешная попытка.
	err = bc.AddShipments(shipments_company2,
		AuthorizedValidators[1])
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Неуспешная попытка.
	err = bc.AddShipments(shipments_company3,
		"bad_validator")
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Вывод всей цепочки.
	out, _ := json.MarshalIndent(bc.Chain, "", " ")
	fmt.Println(string(out))

	// Проверка валидности.
	fmt.Println("Chain valid?", bc.IsValid())
}
