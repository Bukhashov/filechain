package user


type UserModel struct {
	ID 			string `json:"id"`
	Name		string `json:"name"`
	Email 		string `json:"email"`
	Image		string `json:"image"`
}