package model

import (
	"apiserver/pkg/auth"
	"apiserver/pkg/constvar"
	"apiserver/pkg/errno"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type UserModel struct {
	BaseModel
	UserName string `json:"username" gorm:"column username;not null" binding:"required" validate:"min=1;max=32"`
	PassWord string `json:"password" gorm:"column username;not null" binding:"required" validate:"min=1;max=32"`
}

func (u *UserModel) TableName() string{
	return "tb_users"
}

// Create creates a new user account.
func (u *UserModel) Create() error{
	return DB.Self.Create(&u).Error
}
// Update updates an user account information.
func (u *UserModel) Update() error{
	return DB.Self.Save(u).Error
}

// ListUser List all users
func ListUser(username string, offset, limit int) ([]*UserModel, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	users := make([]*UserModel, 0)
	var count uint64

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Self.Model(&UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

// GetUser gets an user by the user identifier.
func GetUser(username string) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("username = ?", username).First(&u)
	return u, d.Error
}


// DeleteUser updates an user account information.
func DeleteUser(id uint64) error{
	user := UserModel{}
	user.BaseModel.Id = id
	return DB.Self.Delete(&user).Error
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *UserModel) Compare(pwd string) (err error) {
	err = auth.Compare(u.PassWord, pwd)
	return
}

// Encrypt user password
func (u *UserModel) Encrypt() (err error){
	u.PassWord, err = auth.Encrypt(u.PassWord)
	if err != nil {
		return err
	}
	return nil
}

// Validate the fields.
func (u *UserModel) Validate() error {
	validate := validator.New()
	return validate.Struct(&u)
}

//checkParams 自定义校验参数.
func (u *UserModel) CheckParams() error{
	if u.UserName == "" {
		return errno.New(errno.ErrValidation,nil).Add("username is empty.")
	}
	if u.PassWord == "" {
		return errno.New(errno.ErrValidation,nil).Add("Password is empty.")
	}
	return nil
}


