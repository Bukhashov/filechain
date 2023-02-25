package user

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/internal/handler/plug"
	"github.com/Bukhashov/filechain/pkg/pb"
	"github.com/Bukhashov/filechain/pkg/utils"
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

func (u *user) Singup(c *gin.Context) {
	// [user] ден сұраудын келген уақыты
	var file multipart.File
	var err error
	
	file, u.Dto.File, err = c.Request.FormFile("img"); if err != nil {
		plug.Response(c, http.StatusBadRequest, "not your img")
		return
	}
	u.Dto.Email = c.PostForm("email")
	u.Dto.Name = c.PostForm("name")

	// Деректер қорына байланыс жасайды
	storage := storage.NewUserStorage(u.client, u.logger);
	// Жаңа [user] жіберген EMAIL басқа [user] ге тиістілі емес екендігін анықтау үшін
	// Деректер қорына сұраныс жібереміз
	// EMAIL басқа [user] ке тиістілі болған жоғдайда
	// Жаңа [user] ден басқа EMAIL қолдануын сұрады
	userModel := model.User{
		Name: u.Dto.Name,
		Email: u.Dto.Email,
	}
	if err = storage.FindUserByEmail(context.TODO(), &userModel); err == nil {
		plug.Response(c, http.StatusBadRequest, "Email address is already in use by another user")
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
		plug.Response(c, http.StatusBadRequest, "Photo should be in [.png .jpg] format only")
		return
	}
	// IMAGE ді сақтау
	tmpFile, err := os.Create(u.Dto.Image); if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}
	defer tmpFile.Close();
	_, err = io.Copy(tmpFile, file); if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
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
		plug.ResponseStatusInternalServerError(c)
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
		plug.ResponseStatusInternalServerError(c)
		return
	}
	if res.Total == 0 || res.Total > 1 {
		plug.Response(c, http.StatusInternalServerError, "There should only be one person in the photo")
		return
	}
	err = storage.Create(context.TODO(), &userModel); if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}
	
	userModel.Image = strconv.FormatInt(userModel.ID, 10)+fileExtension
	err = utils.CopyNewPath(u.Dto.Image, FaceImagePath+userModel.Image); if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}

	err = storage.UpdateIamge(context.TODO(), &userModel); if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}
	plug.Response(c, http.StatusCreated, "created")
}
