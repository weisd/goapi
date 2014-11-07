package models

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/weisd/goapi/modules/log"
)

type Tokens struct {
	Uid  int64 `xorm:NOT NULL`
	Kt   string
	App  string
	Ktkt string
}

var FromType = map[string]string{
	"app":  "app",
	"ktkt": "ktkt",
	"kt":   "kt",
}

/**
 * 生成token
 */
func CreateToken(salt string) string {
	// // 创建uuid
	u1 := uuid.NewV4()
	u5 := uuid.NewV5(u1, salt)
	return u5.String()
}

/**
 *  取token 信息
 */
func GetTokenInfo(uid int64) (*Tokens, error) {
	tokens := &Tokens{Uid: uid}
	has, err := x.Get(tokens)
	if err != nil {
		log.Warn("xorm getTokenInfo failed : %v", err)
		return nil, err
	}
	if !has {
		return nil, errors.New("token info no found")
	}

	return tokens, nil
}

func UpdateToken(uid int64, from string) (token string, err error) {
	salt := fmt.Sprintf("%d", uid)
	token = CreateToken(salt)

	data := make(map[string]interface{})

	key, ok := FromType[from]
	if !ok {
		err = errors.New("from type no found")
		return
	}

	data[key] = token

	// 是否存在
	tokens := &Tokens{Uid: uid}
	has, err := x.Get(tokens)
	if err != nil {
		return
	}
	// 更新
	if has {
		_, err = x.Table(new(Tokens)).Id(uid).Update(data)
		return
	}

	// 创建

	newTokens := new(Tokens)
	newTokens.Uid = uid
	switch from {
	case "ktkt":
		newTokens.Ktkt = token
	case "kt":
		newTokens.Kt = token
	case "app":
		newTokens.App = token
	}

	_, err = x.Insert(newTokens)
	return
}
