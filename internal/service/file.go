package service

import (
	"time"
	"bytes"
	"crypto/sha256"
	"strconv"
	"github.com/Bukhashov/filechain/internal/model"
)

type File struct {
	TimeStamp	int64
	Hash		[]byte
	PrevHash	[]byte
	
	Title		[]byte
	Type		[]byte
	File		[]byte
}

func NewGenesisFile() *model.File {
	genesisFile := NewFile(&model.File{
		TimeStamp: 	time.Now().Unix(),
		PrevHash:	[]byte{},
		Hash: 		[]byte{},
		Title: 		[]byte("Genesis head file"),
		Type: 		[]byte(".run"),
		File: 		[]byte{},
		Access: 	false,
	})
	return genesisFile
}

func NewFile(prevFile *model.File) *model.File {
	file := &File{
		TimeStamp: 	prevFile.TimeStamp,
		Hash: 		[]byte{},
		PrevHash: 	prevFile.PrevHash,
		
		Title: 		prevFile.Title,
		Type: 		prevFile.Type,
		File: 		prevFile.File,
	}
	
	file.setHash()

	newFile := &model.File{
		TimeStamp: 	file.TimeStamp,
		Hash: 		file.Hash,
		PrevHash: 	file.PrevHash,
		
		Title: 		file.Title,
		Type: 		file.Type,
		File: 		file.File,
		Access: 	false,
	}

	return newFile
}

func (f *File) setHash() {
	timeStamp := []byte(strconv.FormatInt(f.TimeStamp, 10))
	headers := bytes.Join([][]byte{f.PrevHash, []byte(f.Title), timeStamp}, []byte{})
	hash := sha256.Sum256(headers)
	f.Hash = hash[:]
}