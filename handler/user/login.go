package user

import (
	. "apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/auth"
	"apiserver/pkg/errno"
	"apiserver/pkg/token"

	"github.com/gin-gonic/gin"
)


// Login 如果密码与指定的帐户匹配 登录生成身份验证令牌
func Login(c *gin.Context) {
	// 将数据与用户结构绑定.
	var u CreateRequest
	if err := c.Bind(&u); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 通过登录用户名获取用户信息.
	d, err := model.GetUser(u.Username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// 比较登录密码和用户密码.
	if err := auth.Compare(d.PassWord, u.Password); err != nil {
		SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// 签署json web令牌
	t, err := token.Sign(c, token.Context{ID: d.Id, Username: d.UserName}, "")
	if err != nil {
		SendResponse(c, errno.ErrToken, nil)
		return
	}

	SendResponse(c, nil, model.Token{Token: t})
}
