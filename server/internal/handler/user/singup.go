package user

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/Bukhashov/filechain/pkg/pb"
	"github.com/Bukhashov/filechain/pkg/utils"
	"github.com/Bukhashov/filechain/internal/model"
)

// FUNC			Жаңа [user] ді тіркеу
// Method:		POST
// ENDPOINT:	[lonstname]:port/api/v1/user/singup
// RESUEST:		form-data
//
// DATA
// img	 type	[]bype
// name  type	string
// email type	string
//
// RESPONSE		
// Content-Type [application/json]
// STATUS 200
// data -> accepted		type	Time
// data -> give away	type	Time
// massage 				type	string
//
// STATUS 400
// data -> accepted		type	Time
// data -> give away	type	Time
// massage 				type	string

func (u *user) Singup(w http.ResponseWriter, r *http.Request){
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
		var file multipart.File
		// [IMG] дін сыйымдылығы < 32Mb дейін 
		// > 32Mb болған жағдайда [status: 400] қайтарылады
		err := r.ParseMultipartForm(32 << 20); if err != nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "Photo size up to 32 mb",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		// #####
		u.Dto.Email = r.FormValue("email")
		u.Dto.Name = r.FormValue("name")
		if u.Dto.Email == "" || u.Dto.Name == ""{
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

		file, u.Dto.File, err = r.FormFile("img"); if err != nil {
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
		
		// Деректер қорына байланыс жасайды
		storage := NewStorage(u.client, u.logger);
		// Жаңа [user] жіберген EMAIL басқа [user] ге тиістілі емес екендігін анықтау үшін
		// Деректер қорына сұраныс жібереміз
		// EMAIL басқа [user] ке тиістілі болған жоғдайда
		// Жаңа [user] ден басқа EMAIL қолдануын сұрады
		err = storage.EmailControl(context.TODO(), &u.Dto); if err == nil {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "Email address is already in use by another user",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}

		// IMAGE дін аты мен форматын бөлкен алу
		// FileName := ImageName
		// Extensiom := .png
		fileNameWithoutExtension, fileExtension := utils.ParseFileName(u.Dto.File.Filename)
		u.Dto.Image = TmpImagePath + u.Dto.File.Filename

		// IMAGE тек [.png .jpg] форматарында қабылданады
		// Жаңа [user] басқа форматтар IMAGE жібермегендігін тексеру
		ok := utils.ControlImgFormat(fileExtension); if !ok {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "Photo should be in [.png .jpg] format only",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		
		// IMAGE ді сақтау
		tmpFile, err := os.Create(u.Dto.Image); if err != nil {
			u.logger.Info(err)
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
		defer tmpFile.Close();

		_, err = io.Copy(tmpFile, file); if err != nil {
			u.logger.Info(err)
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
		
		// протокол gRPC арқылы IMAGE де адамнын бейнесін анықтайтын server сұрау жібереміз
		// IMAGE де қайша адамнын бейнесі бар екендігін анықтайды
		// IMAGE де бір адамнын бейнесі болуы керек
		// Егер бірнеше немесе оданда көп адам бетінін бейнесі болған жағдайда
		// Жаңа [user] ге басқа IMAGE жіберуін және онда тек өзінін бейнесе ғана болуын сұрайды
		stream, err := u.service.Find(context.TODO()); if err != nil {
			u.logger.Info(err)
			return
		}
		// Meta data: FileName және Extension [.png .jpg]
		err = stream.Send(&pb.FindRequest{
			FindData: &pb.FindRequest_Metadata{
				Metadata: &pb.Metadata{
					Filename: fileNameWithoutExtension,
					Extension: fileExtension,
				},
			},
		});
		if err != nil {
			fmt.Printf("find request grpc: %v ", err)
		}
		// Сақталған файлды IMAGE ді ашады
		// Және оқу аяқталған соң жабады
		openFile, err := os.Open(u.Dto.Image); if err != nil {
			u.logger.Info(err)
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
		defer openFile.Close()
		
		reader := bufio.NewReader(openFile)
		// Буфердін ұзындағы 1024 byte
		buffer := make([]byte, 1024)
		
		for {
			n, err := reader.Read(buffer); if err == io.EOF { 
				break 
			}
			if err != nil { 
				u.logger.Info(err) 
			}
			
			err = stream.Send(&pb.FindRequest {
				FindData: &pb.FindRequest_Image {
					Image: buffer[:n],
				},
			});
			if err != nil {
				u.logger.Info(err)
			}
		}
		// протокол gRPC арқылы жіберу аяқталған сон соединения жабылады
		res, err := stream.CloseAndRecv(); if err != nil {
			u.logger.Info(err)
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
	
		if res.Total == 0 {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "There should only be one person in the photo",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		if res.Total > 1 {
			dataResponse := &BadRequrest{
				Data: Data{
					Accepted: timeAccepted,
					GiveAway: time.Now(),
				},
				Massage: "There should only be one person in the photo",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
	
		userModel := &model.User{
			Name: u.Dto.Name, 
			Email: u.Dto.Email,
		}
	
		err = storage.Create(context.TODO(), userModel); if err != nil {
			u.logger.Info(err)
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
	
		userModel.Image = strconv.FormatInt(userModel.ID, 10)+fileExtension
	
		fileManager := utils.NewFileManager()
		err = fileManager.CopyNewPath(u.Dto.Image, FaceImagePath+userModel.Image); if err != nil {
			u.logger.Info(err)
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
	
		err = storage.UpdateIamge(context.TODO(), userModel); if err != nil {
			u.logger.Info(err)
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
		w.WriteHeader(http.StatusCreated)
	}
}