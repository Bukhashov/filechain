package folder

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Bukhashov/filechain/internal/handler/user"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/service"
	"github.com/Bukhashov/filechain/internal/storage"
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
		
		f.Dto.FolderName = r.FormValue("folderName")

		newfolder := service.NewFolder()
		newFile := service.NewGenesisFile()

		f.Model = model.Folder{
			Name: f.Dto.FolderName,
			Addres: newfolder.Addres,
			File: newFile.Hash,
			UserId: userControl.ID,
			// UserId: userControl.ID,
			Access: false,
		}

		// save 
		folderStg := storage.NewFolderStorage(f.client, f.logger)
		err = folderStg.Create(context.TODO(), &f.Model); if err != nil {
			f.logger.Info(err)
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "An error occurred on the server, please try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dataResponse)
			return 
		}

		fileStg := storage.NewFileStorage(f.client, f.logger);
		if err = fileStg.New(context.TODO(), *newFile); err != nil {
			f.logger.Info(err)
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "An error occurred on the server, please try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		// history
		

		dataResponse := &Requrest{
			Data: Data{
				Accepted: timeAccepted,
				GiveAway: time.Now(),
			},
			Addres: string(newfolder.Addres),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dataResponse)

		return
		
	}
}
