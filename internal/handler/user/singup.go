package user

import (
	"github.com/gin-gonic/gin"
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

func (u *user) Singip(c * gin.Context) {

}