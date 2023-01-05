package user

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/Bukhashov/filechain/pkg/pb"
	"github.com/Bukhashov/filechain/pkg/utils"
)
type Singin struct {
	Data	Data
	Massage	string
	Token	string
}

// FUNC			Кіру
// Method:		POST
// ENDPOINT:	[lonstname]:port/api/v1/user/singin
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
// Token				type	string
//
// STATUS 400
// data -> accepted		type	Time
// data -> give away	type	Time
// massage 				type	string

func (u *user) Singin(w http.ResponseWriter, r *http.Request){
	// [user] ден сұраудын келген уақыты
	timeAccepted := time.Now()

	//	Ерер [user] GET әдісмен сұрау жіберген жағдауда бұл әдіс қате және
	// POST әдісмен жіберуін ескертіледі
	if r.Method != http.MethodPost {
		dataResponse := &BadRequrest{
			Data: Data{
				Accepted: timeAccepted,
				GiveAway: time.Now(),
			},
			Massage: "Don't worked method GET use POST methods",
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
				Massage: "Fill in correctly data",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dataResponse)
			return
		}
		// #####
		u.Dto.Email = r.FormValue("email")
		if u.Dto.Email == "" {
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
			
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		storage := NewStorage(u.client, u.logger);
		err = storage.FindByEmail(context.TODO(), u); if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error mail"))
			return
		}
		u.Dto.Image = TmpImagePath + u.Dto.File.Filename
		TmpFileNameWithoutExtension, TmpFileExtension := utils.ParseFileName(u.Dto.File.Filename)
		OriginalFileNameWithoutExtension, OriginalFileExtension := utils.ParseFileName(u.Model.Image)
		
		// Controller format img [.png, .jpg]
		ok := utils.ControlFormat(TmpFileExtension); if !ok {
			fmt.Print("png")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		tmpFile, err := os.Create(u.Dto.Image); if err != nil {
			fmt.Print("creat")
			u.logger.Info(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer tmpFile.Close();
		_, err = io.Copy(tmpFile, file); if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	
		// gRPC
		stream, err := u.service.Comparison(context.TODO()); if err != nil {
			u.logger.Info(err)
			return
		}
		stream.Send(&pb.ComparisonRequest{
			ComparisonData: &pb.ComparisonRequest_OriginalMetadata{
				OriginalMetadata: &pb.Metadata{
					Filename: OriginalFileNameWithoutExtension,
					Extension: OriginalFileExtension,
				},
			},
		})
		
		openTmpFile, err := os.Open(u.Dto.Image); if err != nil{
			u.logger.Info(err)
			return
		}
		defer openTmpFile.Close()
	
		reader := bufio.NewReader(openTmpFile)
		buffer := make([]byte, 1024)
	
		for {
			n, err := reader.Read(buffer); if err == io.EOF {
				break
			}
			if err != nil {
				u.logger.Info(err)
			}
	
			err = stream.Send(&pb.ComparisonRequest{
				ComparisonData: &pb.ComparisonRequest_OriginalImage{
					OriginalImage: buffer[:n],
				},
			})
			if err != nil {
				u.logger.Info(err)
			}
		}
		
		stream.Send(&pb.ComparisonRequest{
			ComparisonData: &pb.ComparisonRequest_ForCheck{
				ForCheck: &pb.Metadata{
					Filename: TmpFileNameWithoutExtension,
					Extension: TmpFileExtension,
				},
			},
		})
		openOriginalFile, err := os.Open(FaceImagePath+u.Model.Image); if err != nil{
			u.logger.Info(err)
			return
		}
		defer openOriginalFile.Close()
		reader = bufio.NewReader(openOriginalFile)
		for {
			n, err := reader.Read(buffer); if err == io.EOF{
				break
			}
			if err != nil {
				u.logger.Info(err)
			}
	
			err = stream.Send(&pb.ComparisonRequest{
				ComparisonData: &pb.ComparisonRequest_ForCheckImage{
					ForCheckImage: buffer[:n],
				},
			})
			if err != nil {
				u.logger.Info(err)
			}
		}
	
		res, err := stream.CloseAndRecv(); if err != nil {
			u.logger.Info(err)
			return
		}
	
		if !res.Coincidences {
			w.WriteHeader(http.StatusBadRequest)
			u.logger.Info(err)
			return
		}
		
		err = u.GeneratorJWT(); if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			u.logger.Info(err)
			return
		}

		// [user] ден келген IMG ді tmp уакытша сақталатын файлдар тізімінен өшіріледі
		if err = os.Remove(u.Dto.Image); err != nil {
			u.logger.Info(err)
		}

		dataSingin := &Singin{
			Data: Data{
				Accepted: timeAccepted,
				GiveAway: time.Now(),
			},
			Massage: "OK",
			Token: u.Token.jwt,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dataSingin)
	}	
}	
