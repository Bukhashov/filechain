package model

type File struct {
	Id			string
	TimeStamp	int64
	Hash 		[]byte
	Type		[]byte
	Title		[]byte
	File 		[]byte
	PrevHash 	[]byte
	Access		bool
}