package history

type HistoryModel struct {
	ID			string		`json:"id"`
	// HEADER
	Hash 		[]byte		`json:"hash"`
	TimeStamp 	int64		`json:"timeStamp"`
	Addres		[]byte		`json:"addres"`
	ProveHash 	[]byte		`json:"poverHash"`
	// BODY
	User		[]byte		`json:"user"`
	Target		[]byte		`json:"target"`
}