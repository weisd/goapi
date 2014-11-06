package models

import (
	"errors"
	"github.com/weisd/goapi/modules/log"
)

type Tokens struct {
	KtToken   string
	AppToken  string
	KtKtToken string
	Uid       int64 `xorm:NOT NULL`
}

func SaveToken() {

}

func CreateToken() {

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
