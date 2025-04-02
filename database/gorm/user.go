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
			// ...其他错误码
		}
		slog.Error("注册失败", "error", err)
		return errors.New("注册失败")
	}
	return nil
}

// 注销用户
func LogOffUser(uid int) error {
	user := model.User{Id: uid}
	tx := globalDB.Delete(&user)
	if tx.Error != nil {
		slog.Error("注销失败", "error", tx.Error)
		return errors.New("注销失败")
	}
	if tx.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

// 获取所有用户
func GetAllUsers() *[]model.User {
	var users []model.User
	err := globalDB.Model(&model.User{}).Find(&users).Error
	if err != nil {
		return nil
	}
	return &users
}

func GetUserByUsername(username string) *model.User {
	user := &model.User{}
	tx := globalDB.Select("*").Where("name = ?", username).First(user)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			slog.Error("获取用户失败", "error", tx.Error)
		}
		return nil
	}
	return user
}

func GetUserByID(uid int) *model.User {
	var user = model.User{Id: uid}
	err := globalDB.Select("*").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("用户不存在", "uid", uid)
		}
		return nil
	}
	return &user
}

func UpdateUserName(uid int, name string) error {
	tx := globalDB.Model(&model.User{}).Where("id = ?", uid).Update("name", name)
	if tx.Error != nil {
		slog.Error("更新用户名失败", "error", tx.Error)
		return errors.New("更新用户名失败")
	} else {
		if tx.RowsAffected <= 0 {
			return errors.New("用户id不存在")
		} else {
			return nil
		}
	}
}

func UpdatePassword(uid int, newPass, oldPass string) error {
	tx := globalDB.Model(&model.User{}).Where("id=? and password=?", uid, oldPass).Update("password", newPass)
	if tx.Error != nil {
		slog.Error("更新密码失败", "error", tx.Error)
		return errors.New("更新密码失败")
	} else {
		if tx.RowsAffected <= 0 {
			return errors.New("旧密码错误")
		} else {
			return nil
		}
	}
}
