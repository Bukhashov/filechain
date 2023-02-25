package folder

import (
	"context"
	"fmt"
	"net/http"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/dto"	
	"github.com/Bukhashov/filechain/pkg/token"
	"github.com/Bukhashov/filechain/internal/handler/plug"
	// "github.com/Bukhashov/filechain/internal/handler/user"

	"github.com/gin-gonic/gin"
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

func (f *folder) GetAllAddress(c *gin.Context) {
	userToken := c.Request.Header["Authorization"]
	
	userDto := dto.User{}
	jwtToken := token.NewToken(f.config.Token.Key)
	err := jwtToken.Parse(userToken[0], &userDto); if err != nil {
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
	
	s := storage.NewFolderStorage(f.client, f.logger)
	ras, err := s.GetAllAddress(context.TODO(), &userDto); if err != nil {
		f.logger.Info(err)
		return
	}
	// ###
	fmt.Print(ras)
}