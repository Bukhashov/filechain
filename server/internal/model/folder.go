package model

type Folder struct {
	ID 			string
	Name 		string
	Addres 		[]byte
	File		[]byte
	UserId		int64
	PrivateKey 	[]byte
	PublicKey	[]byte
	Access 		bool
}