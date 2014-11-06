package account

import (
	"github.com/Unknwon/macaron"
	m "github.com/weisd/goapi/models"
	"github.com/weisd/goapi/modules/log"
	"github.com/weisd/goapi/modules/middleware"
	"strings"
)

func Create(ctx *middleware.Context) {
	uname := ctx.Query("username")
	email := ctx.Query("email")
	password := ctx.Query("password")
	phone := ctx.Query("phone")

	if len(uname) == 0 || len(email) == 0 || len(password) == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	// @todo 验证 类型
	// 用户名不能以数字开头，不能纯数字, 不能有@ ,  防止与email 手机号重复
	log.Debug("start create user ")

	user := &m.User{UserName: uname, Email: email, Password: password, Phone: phone}

	uid, err := m.CreateUser(user)
	if err != nil {
		ctx.ErrorJSON(502, "create user failed")
		return
	}

	ctx.SuccessJSON(uid)
	return
}

/**
 * 验证用户 验证登陆用
 * @params mixed  username|email|phone
 * @params string password
 * @return
 */
func Auth(ctx *middleware.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")

	if len(username) == 0 || len(password) == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	// 全数字，11位 手机验证
	// 有@，email谁
	// username
	var user *m.User
	var err error
	switch {
	case len(username) == 11:
		user, err = m.GetUserByPhone(username)
	case strings.Index(username, "@") > 0:
		user, err = m.GetUserByEmail(username)
	default:
		user, err = m.GetUserByUsername(username)
	}

	if err != nil {
		ctx.ErrorJSON(404, err.Error())
		return
	}

	// 验证成功 创建token 返回

	return
}

func Delete(ctx *macaron.Context) string {
	return "account delete"
}

func Info(ctx *macaron.Context) string {
	return "account delete"
}

func Accounts(ctx *macaron.Context) string {
	return "account delete"
}

func Update(ctx *macaron.Context) string {
	return "account delete"
}
