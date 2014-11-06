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
	Id       int64
	UserName string `xorm:"UNIQUE NOT NULL"`
	Email    string `xorm:"UNIQUE NOT NULL"`
	Phone    string
	Password string    `xorm:"NOT NULL"`
	Salt     string    `xorm:"NOT NULL"`
	Created  time.Time `xorm:"CREATED"`
	Updated  time.Time `xorm:"UPDATED"`
}

// 创建用户
func CreateUser(u *User) (int64, error) {
	u.Salt = GetUserSalt()
	u.EncodePasswd()
	_, err := x.Insert(u)
	if err != nil {
		log.Warn("create user failed : %s", err)
		return 0, err
	}

	return u.Id, nil
}

// EncodePasswd encodes password to safe format.
func (u *User) EncodePasswd() {
	newPasswd := base.PBKDF2([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// GetUserSalt returns a user salt token
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
	user = &User{UserName: username}
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
