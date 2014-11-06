package token

import (
	"github.com/Unknwon/macaron"
	"github.com/satori/go.uuid"
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
	// // 创建uuid
	u1 := uuid.NewV4()
	u5 := uuid.NewV5(u1, salt)

	log.Debug("create token : %s", u5.String())

	ctx.SuccessJSON(u5.String())
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
		authToken = tokenInfo.KtToken
	case "ktkt":
		authToken = tokenInfo.KtKtToken
	case "app":
		authToken = tokenInfo.AppToken
	}

	if authToken != token {
		ctx.ErrorJSON(403, "token auth failed")
		return
	}

	ctx.SuccessJSON("ok")
	return
}

func Delete(ctx *macaron.Context) string {
	return "token delete"
}
