package database

import (
	"errors"
	model "infomation-release/database/model"

	"log/slog"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func RegisterUser(username, password string) error {
	var user = new(model.User) // 一定要创建一块内存，初始化这个结构体
	user.Name = username
	user.PassWord = password
	err := globalDB.Create(user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				slog.Error("用户已存在", "username", username)
				return errors.New("用户已存在")
			}
		}
		slog.Error("注册失败", "error", err)
		return errors.New("注册失败")
	}
	return nil
}

func GetAllUsers() *[]model.User {
	var users []model.User
	err := globalDB.Model(&model.User{}).Find(&users).Error
	if err != nil {
		return nil
	}
	return &users
}

func GetUserByID(uid int) *model.User {
	var user = model.User{Id: uid}
	tx := globalDB.Select("*").First(&user)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("用户不存在", "uid", uid)
		}
		return nil
	}
	return &user
}
