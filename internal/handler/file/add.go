package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/Bukhashov/filechain/internal/handler/plug"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/dto"
	"github.com/Bukhashov/filechain/internal/service"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/pkg/utils"
	"github.com/Bukhashov/filechain/pkg/token"

	"github.com/gin-gonic/gin"
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
// folder_address	type	string
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

func (f *filechain) Add(c *gin.Context){
	userToken := c.Request.Header["Authorization"]

	if len(userToken) < 1 {
		plug.Response(c, http.StatusBadRequest, "not fount token, plass auth")
	}
	
	userDto := dto.User{}
	
	jwtToken := token.NewToken(f.config.Token.Key)
	
	err := jwtToken.Parse(userToken[0][7:], &userDto); if err != nil {
		plug.Response(c, http.StatusBadRequest, "err token")
		return
	}

	userModel := &model.User{
		Email: userDto.Email,
	}
	UserStorage := storage.NewUserStorage(f.client, f.logger);
	err = UserStorage.FindUserByEmail(context.TODO(), userModel); if err != nil {
		plug.Response(c, http.StatusBadRequest, err.Error())
		return
	}
		

	var file multipart.File
	// [FILE] дін сыйымдылығы < 32Mb дейін 
	// > 32Mb болған жағдайда [status: 400] қайтарылады
	// *****
	file, f.Dto.File, err = c.Request.FormFile("file"); if err != nil {
		plug.Response(c, http.StatusBadRequest, "Fill in correctly data")
		return
	}
	f.Dto.Title = c.PostForm("title");
	f.Dto.Address = c.PostForm("folder_address");
	folderStorage := storage.NewFolderStorage(f.client, f.logger)

	// fmt.Printf("user id %v \n file address %v  end\n", userControl.ID, f.Dto.Address)
		
	folder, err := folderStorage.GetFolder(context.TODO(),  &userDto.ID, []byte(f.Dto.Address)); if err != nil {
		plug.Response(c, http.StatusBadRequest, "you don't have such an address file")	
		return
	}
	fmt.Printf("user id %v\nfolder hash: %v", userDto.ID, folder.File)
	
	// FILE дін аты мен форматын бөлкен алу
	// FileName := ImageName
	// Extensiom := .png
	fileNameWithoutExtension, fileExtension := utils.ParseFileName(f.Dto.File.Filename)
	f.Dto.FilePath = TmpFilePath + f.Dto.File.Filename

	// IMAGE тек [.pdf] форматарында қабылданады
	ok := utils.ControlFileFormat(fileExtension); if !ok {
		plug.Response(c, http.StatusBadRequest, "FILE should be in [.pdf] format only")		
		return
	}
	// IMAGE ді сақтау
	
	tmpFile, err := os.Create(f.Dto.FilePath); if err != nil {
		f.logger.Info(err)
		plug.Response(c, http.StatusBadRequest, "create an error occurred on the server, please try again later")	
			return
	}
	defer tmpFile.Close();
		
	_, err = io.Copy(tmpFile, file); if err != nil {
		f.logger.Info(err)
		plug.Response(c, http.StatusBadRequest, "copy An error occurred on the server, please try again later")
		return
	}

	fileByte, err := os.ReadFile(f.Dto.FilePath); if err != nil {
		f.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
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
		UserId: userDto.ID,
		Addres: []byte(f.Dto.Address),
	})

	fileStorage := storage.NewFileStorage(f.client, f.logger)
	if err = fileStorage.New(context.TODO(), *newBlockFile); err != nil {
		f.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}

	plug.Response(c, http.StatusOK, "saved")
}