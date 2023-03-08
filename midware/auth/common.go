package midware

import (
	"github.com/kataras/iris/v12"
	"github.com/lishimeng/app-starter/token"
)

type Response struct {
	Code    interface{} `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
}

func errorWith(ctx iris.Context, code int, err error) {
	var resp Response
	resp.Code = code
	resp.Message = err.Error()
	_ = ctx.JSON(resp)
}

func GetUid(ctx iris.Context) (uid string) {
	uid = ctx.GetHeader(UidKey)
	return
}

func GetClientType(ctx iris.Context) (ct string) {
	ct = ctx.GetHeader(ClientKey)
	return
}

func GetOrg(ctx iris.Context) (org string) {
	org = ctx.GetHeader(OrgKey)
	return
}

func GetDept(ctx iris.Context) (dept string) {
	dept = ctx.GetHeader(DeptKey)
	return
}

func GetUserInfo(ctx iris.Context) (uid token.JwtPayload, err error) {
	var ui = ctx.Values().Get(UserInfoKey)
	uid, ok := ui.(token.JwtPayload)
	if !ok {
		err = ErrNotAllowed
		return
	}
	return
}
