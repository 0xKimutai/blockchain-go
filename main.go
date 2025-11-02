package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type Block struct {
	timestamp  		int64
	data 			[]byte
	prevBlockHash 	[]byte
	hash 			[]byte
}

// Calculate the hash of the block
func (b *Block) setHash() {
	headers := bytes.Join(
		// timestamp, data, prevhash
		[][]byte{
			IntToHex(b.timestamp),
			b.data,
			b.prevBlockHash,
		},
		[]byte{},
	)
	// convert headers to hash
	hash := sha256.Sum256(headers)
	b.hash = hash[:]
}

// create function IntToHex to convert int64 to byte array
func IntToHex( x int64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

// main function
func main() {
	b := &Block{
		timestamp: time.Now().Unix(),
		data: []byte("Hello, blockchain!"),
		prevBlockHash: []byte{},
	}

	b.setHash()
	// %x to print byte array in hex format
	fmt.Printf("Blockhash: %x \n", b.hash)
}