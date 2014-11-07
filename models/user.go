package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/weisd/goapi/modules/base"
	"github.com/weisd/goapi/modules/log"
	"time"
)

type User struct {
	Id       int64     `json:"id"`
	Username string    `json:"username",xorm:"UNIQUE NOT NULL"`
	Email    string    `json:"email",xorm:"UNIQUE NOT NULL"`
	Phone    string    `json:"phone"`
	Password string    `json:"_",xorm:"NOT NULL"`
	Salt     string    `json:"_",xorm:"NOT NULL"`
	Created  time.Time `json:"created",xorm:"CREATED"`
	Updated  time.Time `json:"updated",xorm:"UPDATED"`
}

// 创建用户
func CreateUser(u *User) (int64, error) {
	u.Salt = GetUserSalt()
	u.Password = EncodePasswd(u.Password, u.Salt)
	_, err := x.Insert(u)
	if err != nil {
		log.Warn("create user failed : %s", err)
		return 0, err
	}

	return u.Id, nil
}

// 用salt加密password
func EncodePasswd(password, salt string) string {
	newPasswd := base.PBKDF2([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return fmt.Sprintf("%x", newPasswd)
}

// 生成salt
func GetUserSalt() string {
	return base.GetRandomString(10)
}

/**
 * 取用户信息
 * @params int64 uid
 * @return User, error
 */
func GetUserInfo(uid int64) (user *User, err error) {
	user = &User{Id: uid}
	has, err := x.Get(user)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("user no found")
		return
	}

	return
}

func GetUserByPhone(phone string) (user *User, err error) {
	user = &User{Phone: phone}
	has, err := x.Get(user)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("user no found")
		return
	}

	return
}

func GetUserByEmail(email string) (user *User, err error) {
	user = &User{Email: email}
	has, err := x.Get(user)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("user no found")
		return
	}

	return
}

func GetUserByUsername(username string) (user *User, err error) {
	user = &User{Username: username}
	has, err := x.Get(user)
	if err != nil {
		return
	}

	if !has {
		err = errors.New("user no found")
		return
	}

	return
}

// 更新用户信息
// @params int64 uid
// @params map[string]interface{}
// @return error
func UpdateUserInfo(uid int64, data map[string]interface{}) error {
	// 更新密码重新
	if _, ok := data["salt"]; ok {
		delete(data, "salt")
	}
	if _, ok := data["id"]; ok {
		delete(data, "id")
	}
	if pwd, ok := data["password"]; ok {
		data["salt"] = GetUserSalt()
		data["password"] = EncodePasswd(pwd.(string), data["salt"].(string))
	}

	if len(data) == 0 {
		return errors.New("update data is empty")
	}

	_, err := x.Table(new(User)).Id(uid).Update(data)
	return err
}

/**
 * 取用户列表
 */
func GetUserList() (list []User, err error) {

	list = make([]User, 0)

	err = x.Find(&list)
	if err != nil {
		return
	}
	return
}
