package user

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/Bukhashov/filechain/pkg/pb"
	"github.com/Bukhashov/filechain/pkg/utils"
)

func (u *user) Singup(w http.ResponseWriter, r *http.Request){
	var file multipart.File

	err := r.ParseMultipartForm(32 << 20); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// #####
	u.Dto.Email = r.FormValue("email")
	u.Dto.Name = r.FormValue("name")
	if u.Dto.Email == "" || u.Dto.Name == ""{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, u.Dto.File, err = r.FormFile("img"); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	fileNameWithoutExtension, fileExtension := utils.ParseFileName(u.Dto.File.Filename)
	
	// Controller format img [.png, .jpg]
	ok := utils.ControlFormat(fileExtension); if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	tmpFile, err := os.Create("./assets/tmp/" + u.Dto.File.Filename); if err != nil {
		fmt.Print("tmp file err")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer tmpFile.Close();

	_, err = io.Copy(tmpFile, file); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	// gRPC
	stream, err := u.service.Find(context.TODO()); if err != nil {
		fmt.Print(err)
	}
	// send meta data
	err = stream.Send(&pb.FindRequest{
		FindData: &pb.FindRequest_Metadata{
			Metadata: &pb.Metadata{
				Filename: fileNameWithoutExtension,
				Extension: fileExtension,
			},
		},
	}); 
	if err != nil {
		fmt.Print(err)
	}

	openFile, err := os.Open("./assets/tmp/"+u.Dto.File.Filename); if err != nil {
		u.logger.Info(err)
		return
	}

	reader := bufio.NewReader(openFile)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer); if err == io.EOF { 
			u.logger.Info(err)
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
			fmt.Print(err)
		}
	}
	res, err := stream.CloseAndRecv(); if err != nil {
		fmt.Print(err)
		return
	}

	if res.Total == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if res.Total > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	

}