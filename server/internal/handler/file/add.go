package file

import (
	"context"
	"encoding/json"
	"fmt"

	// "fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/Bukhashov/filechain/internal/handler/user"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/service"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/pkg/utils"
)

// FUNC			Жаңа [user] ді тіркеу
// Method:		POST
// ENDPOINT:	[lonstname]:port/api/v1/new/floder
// RESUEST:		form-data
// ------------------------------------
// == HEADER ==
// Authentication	type	string
// == BODY ==
// title  			type	string
// file				type	[]byte
// address			type	string
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

func (f *filechain) Add(w http.ResponseWriter, r *http.Request){
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

		var file multipart.File
		// [FILE] дін сыйымдылығы < 32Mb дейін 
		// > 32Mb болған жағдайда [status: 400] қайтарылады
		err = r.ParseMultipartForm(32 << 20); if err != nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "File size up to 32 mb",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		
		f.Dto.Title = r.FormValue("title"); if f.Dto.Title == "" {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "specify a title",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		f.Dto.Address = r.FormValue("address"); if f.Dto.Title == "" {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "send folder address",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		
		folderStorage := storage.NewFolderStorage(f.client, f.logger)
		

		// fmt.Printf("user id %v \n file address %v  end\n", userControl.ID, f.Dto.Address)
		
		folder, err := folderStorage.GetFolder(context.TODO(),  &userControl.ID, []byte(f.Dto.Address)); if err != nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "you don't have such an address file",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		fmt.Printf("user id %v\nfolder hash: %v", userControl.ID, folder.File)

		file, f.Dto.File, err = r.FormFile("file"); if err != nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "Fill in correctly data",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		// FILE дін аты мен форматын бөлкен алу
		// FileName := ImageName
		// Extensiom := .png
		fileNameWithoutExtension, fileExtension := utils.ParseFileName(f.Dto.File.Filename)
		f.Dto.FilePath = TmpFilePath + f.Dto.File.Filename

		// IMAGE тек [.pdf] форматарында қабылданады
		ok := utils.ControlFileFormat(fileExtension); if !ok {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "FILE should be in [.pdf] format only",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		// IMAGE ді сақтау
		tmpFile, err := os.Create(f.Dto.FilePath); if err != nil {
			f.logger.Info(err)
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "create an error occurred on the server, please try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		defer tmpFile.Close();
		
		_, err = io.Copy(tmpFile, file); if err != nil {
			f.logger.Info(err)
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "copy An error occurred on the server, please try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		fileByte, err := os.ReadFile(f.Dto.FilePath); if err != nil {
			f.logger.Info(err)
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "read file An error occurred on the server, please try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		// fmt.Printf("%v ", &folder.File)

		newBlockFile := service.NewFile(&model.File{
			TimeStamp:	time.Now().Unix(),
			PrevHash: 	folder.File,
			Type: 		[]byte(fileExtension),
			Title: 		[]byte(f.Dto.Title),
			FileName: 	[]byte(fileNameWithoutExtension),
			File:		fileByte,
			Access: 	false,
		})

		folderStorage.Update(context.TODO(), &model.Folder{
			File: newBlockFile.Hash,
			UserId: userControl.ID,
			Addres: []byte(f.Dto.Address),
		})

		fileStorage := storage.NewFileStorage(f.client, f.logger)
		if err = fileStorage.New(context.TODO(), *newBlockFile); err != nil {
			f.logger.Info(err)
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "new file storage An error occurred on the server, please try again later",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		dataResponse := &Requrest{
			Data: Data{
				Accepted: timeAccepted,
				GiveAway: time.Now(),
			},
			Massage: "saved",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dataResponse)
		return
	}
}