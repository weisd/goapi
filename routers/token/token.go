package token

import (
	m "github.com/weisd/goapi/models"
	"github.com/weisd/goapi/modules/log"
	"github.com/weisd/goapi/modules/middleware"
)

type returnErr struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// 创建token
func Create(ctx *middleware.Context) {

	salt := ctx.Query("salt")

	if len(salt) == 0 {
		salt = "ktkt"
	}
	token := m.CreateToken(salt)

	ctx.SuccessJSON(token)
}

/**
 * 更新用户token
 * @params int64 uid
 * @params string from
 * @return string token
 */
func Update(ctx *middleware.Context) {
	uid := ctx.QueryInt64("uid")
	from := ctx.Query("from")

	if uid == 0 || len(from) == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	if _, ok := m.FromType[from]; !ok {
		ctx.ErrorJSON(401, "from type error")
		return
	}

	token, err := m.UpdateToken(uid, from)
	if err != nil {
		ctx.ErrorJSON(502, "update token failed :"+err.Error())
		return
	}

	ctx.SuccessJSON(token)
}

// 验证token
// @params string token
// @params string uid
func Auth(ctx *middleware.Context) {
	uid := ctx.QueryInt64("uid")
	token := ctx.Query("token")
	from := ctx.Query("from")

	if uid == 0 || len(token) == 0 || len(from) == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	if _, ok := m.FromType[from]; !ok {
		ctx.ErrorJSON(401, "from type error")
		return
	}

	// 通过uid查数据库
	tokenInfo, err := m.GetTokenInfo(uid)
	log.Debug("%d, %s, %v", uid, token, tokenInfo)

	if err != nil {
		ctx.ErrorJSON(404, "token no found")
		return
	}

	authToken := ""

	switch from {
	case "kt":
		authToken = tokenInfo.Kt
	case "ktkt":
		authToken = tokenInfo.Ktkt
	case "app":
		authToken = tokenInfo.App
	}

	if authToken != token {
		ctx.ErrorJSON(403, "token auth failed")
		return
	}

	ctx.SuccessJSON("ok")
	return
}
