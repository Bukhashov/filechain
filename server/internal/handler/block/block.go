package block

import(
	"net/http"
	"github.com/Bukhashov/filechain/internal/config"
	"github.com/Bukhashov/filechain/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Filechain interface {
	New(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	All(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
}

type filechain struct {
	client 	*pgxpool.Pool
	logger	*logging.Logger
	config	config.Config
}


func NewBlock(client *pgxpool.Pool, config config.Config, logger *logging.Logger) Filechain{
	return &filechain {
		config: config,
		logger: logger,
		client: client,
	}
}











// import (
// 	"bytes"
// 	"crypto/sha256"
// 	"strconv"
// 	"time"
// )

// type Filechain struct {
//     Blocks []*Block
// }

// func(b *Block) setHash(){
// 	timeStamp := []byte(strconv.FormatInt(b.TimeStamp, 10))
// 	headers := bytes.Join([][]byte{b.PrevHash, b.Data, timeStamp}, []byte{})
// 	hash := sha256.Sum256(headers)

// 	b.Hash = hash[:]
// }

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