package plug

import (
	"github.com/gin-gonic/gin"

)

func Response(c *gin.Context, status int, massage string) {
	c.JSON(status, gin.H{
		"code" : status,
		"massage" : massage,
	})
}