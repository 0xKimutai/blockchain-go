package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

type Block struct {
	Timestamp  		int64
	Data 			[]byte
	PrevBlockHash 	[]byte
	Hash 			[]byte
	Nonce			int
}

type Blockchain struct {
	blocks []*Block
}

// Calculate the hash of the block
func (b *Block) SetHash() {
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
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	isValid := pow.Validate()
	fmt.Printf("Block mined: %s\n", strconv.FormatBool(isValid))

	return block
}

//add blocks to blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]	// gets the most recent block in chain
	newBlock := NewBlock(data, prevBlock.Hash)	// create the new block -> link to prev block hash
	bc.blocks = append(bc.blocks, newBlock)		// add the new block to the chain
}

// first block in the chain -> genesis block -> no prev block hash
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// funct to create a blockchain with the genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// PROOF OF WORK
const targetBits = 16

type ProofOfWork struct {
	block *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// data to be hashed
func (pow *ProofOfWork) prepareData (nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// core of PoW algorithm
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining block contains: %s\n", pow.block.Data)

	maxNounce := math.MaxInt64
	for nonce < maxNounce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")

	return nonce, hash[:]
}

// validate proof or works
func (pow *ProofOfWork) Validate () bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	hash :=sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid

}

// main function
func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 0.1 BTC to jon")
	bc.AddBlock("Send 0.2 BTC to mike")

	for _, block := range bc.blocks {
		fmt.Printf("Prev block hash", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format(time.RFC3339))
		fmt.Println()
	}

}