package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

type Block struct {
	timestamp     int64
	data          []byte
	prevBlockHash []byte
	hash          []byte
}

// SetHash calculates and sets the hash of the block.
func (b *Block) SetHash() {
	headers := bytes.Join(
		[][]byte{
			IntToHex(b.timestamp),
			b.data,
			b.prevBlockHash,
		},
		[]byte{},
	)
	hash := sha256.Sum256(headers)
	b.hash = hash[:]
}

// IntToHex converts an int64 to a byte slice.
func IntToHex(n int64) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()
}

// Example usage
func main() {
	b := &Block{
		timestamp:     time.Now().Unix(),
		data:          []byte("Hello, blockchain"),
		prevBlockHash: []byte{},
	}
	b.SetHash()
	fmt.Printf("Block hash: %x\n", b.hash)
}