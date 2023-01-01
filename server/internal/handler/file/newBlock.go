package block

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	// "time"
)

type fileBlock struct {
	Block	[]*Block	
}

func(b *Block) setHash(){
	timeStamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
	headers := bytes.Join([][]byte{b.PrevHash, b.Data, timeStamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// func (filechain *Filechain) AddBlock(data string) {
// 	prevBlock := filechain.Blocks[len(filechain.Blocks)-1]
// 	newBlock := newBlock(data, prevBlock.Hash);
// 	filechain.Blocks = append(filechain.Blocks, newBlock)
// }

// func newGenesisBlock() *Block {
// 	return newBlock("Genesis block", []byte{})
// }

// func newBlock(data string, prevBlockHash []byte) *Block {
// 	block := &Block{
// 		time.Now().Unix(),
// 		[]byte{},
// 		[]byte(data),
// 		prevBlockHash,
// 	}
// 	block.setHash()
	
// 	return block
// }

// func NewFilechain() * Filechain {
// 	return &Filechain{
// 		[]*Block{newGenesisBlock()},
// 	}
// }
