package folder

import (
	"net/http"
	"time"
	"encoding/json"
)

// FUNC			Жаңа [user] ді тіркеу
// Method:		POST
// ENDPOINT:	[lonstname]:port/api/v1/new/floder
// RESUEST:		form-data
// ------------------------------------
// == HEADER ==
// Authentication	type	string
// == BODY ==
// name  			type	string
// ------------------------------------
// RESPONSE		
// Content-Type [application/json]
// STATUS 200
// data -> accepted		type	Time
// data -> give away	type	Time
// massage 				type	string
// addres				type	string
// STATUS 400
// data -> accepted		type	Time
// data -> give away	type	Time
// massage 				type	string

func (f *folder) New(w http.ResponseWriter, r *http.Request) {
	// [user] ден сұраудын келген уақыты
	timeAccepted := time.Now()
	//	Ерер [user] GET әдісмен сұрау жіберген жағдауда бұл әдіс қате және
	//	POST әдісмен жіберуін ескертіледі
	if r.Method != http.MethodPost {
		dataResponse := &BadRequrest{
			Data: Data{
				Accepted: timeAccepted,
				GiveAway: time.Now(),
			},
			Massage: r.Method + "method don't worked use POST methods",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dataResponse)
		return
	}

	if r.Method == http.MethodPost {
		
	}
}