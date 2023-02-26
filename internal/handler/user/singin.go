package user

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"github.com/Bukhashov/filechain/internal/storage"
	"github.com/Bukhashov/filechain/internal/handler/plug"
	"github.com/Bukhashov/filechain/internal/model"
	"github.com/Bukhashov/filechain/pkg/token"
	"github.com/Bukhashov/filechain/pkg/pb"
	"github.com/Bukhashov/filechain/pkg/utils"
	
	"github.com/gin-gonic/gin"
)

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

func (u *user) Singin(c *gin.Context) {
	var file multipart.File
	var err error
	
	u.Dto.Email = c.PostForm("email")
	file, u.Dto.File, err = c.Request.FormFile("img"); if err != nil {
		plug.Response(c, http.StatusBadRequest, "not your img")
		return
	}

	storage := storage.NewUserStorage(u.client, u.logger);
	// userModel := model.User{
	// 	Email: u.Dto.Email,
	// }
	u.Model.Email = u.Dto.Email

	err = storage.FindUserByEmail(context.TODO(), &u.Model); if err != nil {
		plug.Response(c, http.StatusBadRequest, "error mail")
		return
	}
	u.Dto.Image = TmpImagePath + u.Dto.File.Filename
	TmpFileNameWithoutExtension, TmpFileExtension := utils.ParseFileName(u.Dto.File.Filename)
	OriginalFileNameWithoutExtension, OriginalFileExtension := utils.ParseFileName(u.Model.Image)
	// Controller format img [.png, .jpg]
	ok := utils.ControlImgFormat(TmpFileExtension); if !ok {
		plug.Response(c, http.StatusBadRequest, "Photo should be in [.png .jpg] format only")
		return
	}
	tmpFile, err := os.Create(u.Dto.Image); if err != nil {
		fmt.Print("creat")
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}
	defer tmpFile.Close();
	_, err = io.Copy(tmpFile, file); if err != nil {
		plug.ResponseStatusInternalServerError(c)
		return
	}

	// gRPC
	stream, err := u.service.Comparison(context.TODO()); if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
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
		plug.ResponseStatusInternalServerError(c)
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
		plug.ResponseStatusInternalServerError(c)
		return
	}

	if !res.Coincidences {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}

	jwtToken := token.NewToken(u.config.Token.Key)
	newToken, err := jwtToken.Generator(&model.User{
		ID: u.Dto.ID,
		Email: u.Dto.Email,
		Name: u.Dto.Name,
	})
	if err != nil {
		u.logger.Info(err)
		plug.ResponseStatusInternalServerError(c)
		return
	}
	
	// [user] ден келген IMG ді tmp уакытша сақталатын файлдар тізімінен өшіріледі
	if err = os.Remove(u.Dto.Image); err != nil {
		u.logger.Info(err)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code" : http.StatusOK,
		"maggase" : "ok",
		"token" : newToken,
	})

}