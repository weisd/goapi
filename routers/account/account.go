package account

import (
	"encoding/json"
	"fmt"
	m "github.com/weisd/goapi/models"
	"github.com/weisd/goapi/modules/base"
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

	user := &m.User{Username: uname, Email: email, Password: password, Phone: phone}

	uid, err := m.CreateUser(user)
	if err != nil || uid == 0 {
		ctx.ErrorJSON(502, "create user failed")
		return
	}

	ctx.SuccessJSON(user)
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
	from := ctx.Query("from")

	if len(username) == 0 || len(password) == 0 || len(from) == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	// 全数字，11位 手机验证
	// 有@，email谁
	// username
	var user *m.User
	var err error
	switch {
	case base.IsNumber(username) && len(username) == 11:
		user, err = m.GetUserByPhone(username)
	case strings.Index(username, "@") > 0:
		user, err = m.GetUserByEmail(username)
	default:
		user, err = m.GetUserByUsername(username)
	}

	if err != nil || user == nil {
		ctx.ErrorJSON(404, err.Error())
		return
	}

	if m.EncodePasswd(password, user.Salt) != user.Password {
		ctx.ErrorJSON(403, "auth failed")
		return
	}

	// 验证成功 创建token 返回
	token := m.CreateToken(fmt.Sprintf("%d", user.Id))
	log.Debug("user info : %v", token)

	ctx.SuccessJSON(user)

	return
}

/**
 * 取用户信息
 * @params int64 uid
 * @return *User
 */
func Info(ctx *middleware.Context) {
	uid := ctx.QueryInt64("uid")
	if uid == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	user, err := m.GetUserInfo(uid)
	if err != nil || user == nil {
		ctx.ErrorJSON(404, "user no found")
		return
	}

	ctx.SuccessJSON(user)
}

/**
 * 更新用户信息
 * @params int64 uid
 * @params json  data
 * @return bool ok
 */
func Update(ctx *middleware.Context) {
	uid := ctx.QueryInt64("uid")
	data := ctx.Query("data")

	if uid == 0 || len(data) == 0 {
		ctx.ErrorJSON(401, "params error")
		return
	}

	log.Debug("update id : %d, data: %s", uid, data)

	// 定义类型，接收
	var update map[string]interface{}

	err := json.Unmarshal([]byte(data), &update)
	if err != nil {
		ctx.ErrorJSON(401, "json unmarshal failed"+err.Error())
		return
	}

	log.Debug("udpate info map : %v", update)

	err = m.UpdateUserInfo(uid, update)
	if err != nil {
		ctx.ErrorJSON(502, "update failed : "+err.Error())
		return
	}

	ctx.SuccessJSON("ok")
}

/**
 * 取用户列表
 */
func Accounts(ctx *middleware.Context) {
	list, err := m.GetUserList()
	if err != nil {
		ctx.ErrorJSON(502, "GetUserList failed")
		return
	}

	ctx.SuccessJSON(list)
	return
}

func Delete(ctx *middleware.Context) string {
	return "account delete"
}
