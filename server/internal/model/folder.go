package model

type Folder struct {
	ID 			int64
	Name 		string
	Addres 		[]byte
	File		[]byte
	UserId		int64
	PrivateKey 	[]byte
	PublicKey	[]byte
	Access 		bool
}