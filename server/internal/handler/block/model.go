package block

type Block struct {
	TimeStamp	int64
	Hash 		[]byte
	Type		string
	Data 		[]byte
	PrevHash 	[]byte
}