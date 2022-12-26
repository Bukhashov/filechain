package user


type UserModel struct {
	ID 			string `json:"id"`
	Name		string `json:"name"`
	Email 		string `json:"email"`
	Confirm 	bool `json:"confirm"`
	ConfirmCode string `json:"confirm_code"`
}