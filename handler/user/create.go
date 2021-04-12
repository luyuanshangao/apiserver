package user

import (
	. "apiserver/handler"
	"apiserver/model"
	"apiserver/pkg/errno"
	"github.com/gin-gonic/gin"
)
type createResponse struct {
	Username string `json:"username"`
}
type  createRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
// Create 创建新用户帐户.
func Create(c *gin.Context) {

	var r  model.UserModel

	if err := c.BindJSON(&r); err != nil {
		SendResponse(c,errno.ErrBind,gin.H{"error": "not params"})
		return
	}
	u := model.UserModel{
		UserName: r.UserName,
		PassWord: r.PassWord,
	}
	//Validate校验参数
	if err := u.Validate(); err != nil{
		SendResponse(c,err,nil)
		return
	}

	//加密用户密码
	if err := u.Encrypt(); err != nil {
		SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	//将用户插入数据库
	if err := u.Create(); err != nil {
		SendResponse(c, errno.ErrDatabase,nil)
		return
	}

	//自定义校验参数
	if err := u.CheckParams(); err != nil{
		SendResponse(c,err,nil)
		return
	}


	rsp := createResponse{
		Username: r.UserName,
	}

	SendResponse(c,nil,rsp)
}




//.Add("This is add message.") 对外展示更多的信息
// 第一个参数为errno code码 第二个参数为 后台日志可看的敏感信息 用户不会看到
//err = errno.New(errno.ErrPasswordIncorrect, fmt.Errorf("password can not found in db: xx.xx.xx.xx"))
//handler.SendResponse(c,err,nil)
//log.Debugf("username is: [%s]", r.Username) //添加后台debug日志
//log.Errorf(err, "password is empty") //添加后台error日志 第二个参数为标题 具体错误在err中
//err = fmt.Errorf("password is empty")   //使用默认错误 该错误信息不是定制的错误类型 解析时会解析为默认的 errno.InternalServerError 错误

