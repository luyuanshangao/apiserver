package user

import (
	. "apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/errno"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Update(c *gin.Context)  {

	userId, _ := strconv.Atoi(c.Param("id"))

	var u model.UserModel

	if err := c.Bind(&u); err  != nil {
		SendResponse(c,err,nil)
	}

	// We update the record based on the user id.
	u.Id = uint64(userId)
	// Validate the data.
	if err := u.Validate(); err != nil {
		SendResponse(c, errno.ErrValidation, nil)
		return
	}
	// Encrypt the user password.
	if err := u.Encrypt(); err != nil {
		SendResponse(c, errno.ErrEncrypt, nil)
		return
	}
	// Save changed fields.
	if err := u.Update(); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	SendResponse(c,nil,nil)


}