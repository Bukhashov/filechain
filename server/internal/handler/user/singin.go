package user

import (
	"net/http"
)

func (u *user) Singin(w http.ResponseWriter, req *http.Request){
	u.logger.Info("end point user login")
}