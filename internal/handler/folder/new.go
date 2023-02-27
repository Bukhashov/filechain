package folder

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Bukhashov/filechain/internal/dto"
	// "github.com/Bukhashov/filechain/internal/handler/user"
	"github.com/Bukhashov/filechain/internal/handler/plug"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/service"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/pkg/token"
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


func (f *folder) New(c *gin.Context) {

	userToken := c.Request.Header["Authorization"]
	
	userDto := dto.User{}
		
	jwtToken := token.NewToken(f.config.Token.Key)
		
	err := jwtToken.Parse(userToken[0][7:], &userDto); if err != nil {
		plug.Response(c, http.StatusBadRequest, "err token")
		return
	}

	UserStorage := storage.NewUserStorage(f.client, f.logger);
		
	userModel := &model.User{
		Email: userDto.Email,
	}
	err = UserStorage.FindUserByEmail(context.TODO(), userModel); if err != nil {
		plug.Response(c, http.StatusBadRequest, err.Error())
		return
	}

	f.Dto.FolderName = c.PostForm("folder_name")
		
	newfolder := service.NewFolder()
	newFile := service.NewGenesisFile()

	f.Model = model.Folder{
		Name: f.Dto.FolderName,
		Addres: newfolder.Addres,
		File: newFile.Hash,
		UserId: userDto.ID,
		// UserId: userControl.ID,
		Access: false,
	}

		// save 
		folderStg := storage.NewFolderStorage(f.client, f.logger)
		err = folderStg.Create(context.TODO(), &f.Model); if err != nil {
			f.logger.Info(err)
			plug.ResponseStatusInternalServerError(c)
			return 
		}

		fileStg := storage.NewFileStorage(f.client, f.logger);
		if err = fileStg.New(context.TODO(), *newFile); err != nil {
			f.logger.Info(err)
			plug.ResponseStatusInternalServerError(c)
			return
		}

		// history


		c.JSON(http.StatusOK, gin.H{
			"code" : http.StatusOK,
			"filder_addres" : string(newfolder.Addres),
 		})
		
		return
}