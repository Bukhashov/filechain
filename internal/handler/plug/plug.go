package plug

import (
	"net/http"

	"github.com/Bukhashov/filechain/internal/handler/res"
	"github.com/gin-gonic/gin"
)

func ResponseStatusInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"massage" : "An error occurred on the server, please try again later",
	})
}

func Response(c *gin.Context, status int, massage string) {
	c.JSON(status, gin.H{
		"massage" : massage,
	})
}

func ResponseOk(c *gin.Context, status int, data *res.MsgUserSinginOk) {
	c.JSON(status, data)
}