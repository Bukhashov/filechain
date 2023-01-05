package folder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/internal/handler/user"
)

// FUNC			Жаңа [user] ді тіркеу
// Method:		POST
// ENDPOINT:	[lonstname]:port/api/v1/get/address
// RESUEST:		form-data
// ------------------------------------
// == HEADER ==
// Authentication	type	string
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

func (f *folder) GetAllAddress(w http.ResponseWriter, r *http.Request) {
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
		token := r.Header.Get("Authorization")

		userControl := user.Dto{}
		
		err := userControl.ParseJwt(token); if err != nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "err token",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		if err = userControl.ControlJwt(context.TODO(), f.client, f.logger); err != nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		s := storage.NewFolderStorage(f.client, f.logger)
		ras, err := s.GetAllAddress(context.TODO(), &userControl); if err != nil {
			f.logger.Info(err)
			return
		}
		// ###
		fmt.Print(ras)
	}
}