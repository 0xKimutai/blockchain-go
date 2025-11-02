package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type Block struct {
	Timestamp  		int64
	Data 			[]byte
	PrevBlockHash 	[]byte
	Hash 			[]byte
}

type Blockchain struct {
	blocks []*Block
}

// Calculate the hash of the block
func (b *Block) setHash() {
	headers := bytes.Join(
		// timestamp, data, prevhash
		[][]byte{
			IntToHex(b.Timestamp),
			b.Data,
			b.PrevBlockHash,
		},
		[]byte{},
	)
	// convert headers to hash
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// create function IntToHex to convert int64 to byte array
func IntToHex( x int64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

// function to create a new block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}}
	block.setHash()
	return block
}

//add blocks to blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// first block in the chain -> genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// funct to create a blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// main function
func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 0.1 BTC to jon")
	bc.AddBlock("Send 0.2 BTC to mike")

	for _, block := range bc.blocks {
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("hash:: %x\n", block.Hash)
		fmt.Println()
	}

}