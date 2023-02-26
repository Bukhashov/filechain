package validator

import (
	"fmt"

	"github.com/Bukhashov/filechain/internal/model"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func UserVaidator(u *model.User)(err error){
	err = validation.ValidateStruct(u,
		validation.Field(u.Email, validation.Required, is.Email),
		validation.Field(u.Name, validation.Required, validation.Length(3, 75)),
	)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}