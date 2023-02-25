package res

type MsgStatusInternalServerError struct {
	Massage 	string	`json:"massage"`
}

type MsgUserSinginOk struct {
	Token string `json:"token"`
}