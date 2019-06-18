package member

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"

	"git.yitum.com/saas/shop-admin/pkg/bootstrap"
	"git.yitum.com/saas/shop-admin/router/mdw"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
	"git.yitum.com/saas/shop-admin/model/trans"
	"git.yitum.com/saas/shop-admin/pkg/cache"
	"git.yitum.com/saas/shop-admin/pkg/util"
	"git.yitum.com/saas/shop-admin/router/api"
	"git.yitum.com/saas/shop-admin/router/mdw/auth"
	"git.yitum.com/saas/shop-admin/service"
)

var emailRgx = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-z]{2,4}$`)

// {status: "error", type: "account", currentAuthority: "guest"}
func Login(c *gin.Context) {
	// 如果已经登录
	respView := trans.RespOauthLogin{
		CurrentAuthority: "admin",
	}

	a := auth.Default(c)
	if a.IsAuthenticated() {
		api.JSONOK(c, respView)
		return
	}

	reqView := &trans.ReqOauthLogin{}
	err := c.Bind(reqView)
	if err != nil {
		model.Logger.Info(err.Error(), zap.Int("code", api.MsgErr))
		api.JSON(c, api.MsgErr, "request params is error")
		return
	}

	// 对Identity进行校验，先判断是否是邮箱，若不是邮箱则当做用户名
	isEmail := emailRgx.MatchString(reqView.Identity)
	var biz *mysql.Biz
	if isEmail {
		biz, err = service.Biz.GetBizByPwd("", reqView.Identity, reqView.Pwd, c.ClientIP())
		if err != nil {
			model.Logger.Info("user not record", zap.Int("code", api.MsgErr), zap.String("err", err.Error()))
			api.JSON(c, api.MsgErr, "error")
			return
		}
	} else {
		biz, err = service.Biz.GetBizByPwd(reqView.Identity, "", reqView.Pwd, c.ClientIP())
		if err != nil {
			model.Logger.Info("user not record", zap.Int("code", api.MsgErr), zap.String("err", err.Error()))
			api.JSON(c, api.MsgErr, "error")
			return
		}
	}
	if biz.OpenId == 0 {
		api.JSON(c, api.MsgErr, "error")
	}

	fmt.Println("redirect_uri", reqView.Params.RedirectUri)
	if reqView.Params.RedirectUri == "" {
		session := sessions.Default(c)
		err = auth.AuthenticateSession(session, &auth.Auth{
			Id:       biz.OpenId,
			Nickname: biz.Nickname,
		})
		if err != nil {
			model.Logger.Info(err.Error(), zap.Int("code", api.MsgErr))
			api.JSON(c, api.MsgErr, "error")
			return
		}

		api.JSONOK(c, respView)
		return
	}
}

func Self(c *gin.Context) {
	resp, err := service.Biz.Info(c, mdw.OpenId(c))
	if err != nil {
		model.Logger.Info(err.Error(), zap.Int("code", api.MsgErr), zap.String("err", err.Error()))
		api.JSON(c, api.MsgErr, "fetch account user info error")
		return
	}

	api.JSONOK(c, resp)
	return
}

func SelfPassword(c *gin.Context) {
	req := trans.ReqMemberPassword{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request list params is error")
		return
	}

	_, err := service.Biz.GetBizByOidPwd(auth.Default(c).UniqueId(), req.OldPassword)
	if err != nil {
		model.Logger.Info(err.Error(), zap.Int("code", api.MsgErr), zap.String("err", err.Error()))
		api.JSON(c, api.MsgErr, "password is error")
		return
	}

	err = service.Biz.Update(c, model.Db, mdw.OpenId(c), mysql.Ups{
		"password": req.Password,
	})
	if err != nil {
		model.Logger.Info(err.Error(), zap.Int("code", api.MsgErr), zap.String("err", err.Error()))
		api.JSON(c, api.MsgErr, "password is error")
		return
	}

	api.JSONOK(c)
	return

}

func Logout(c *gin.Context) {
	fmt.Println(auth.Default(c).UniqueId())
	auth.Logout(sessions.Default(c), auth.Default(c))
	api.JSONOK(c)
	return

}

func Add(c *gin.Context) {
	req := &trans.ReqAccountAdd{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request add member params is error")
		return
	}
	if err := service.Biz.Create(c, model.Db, &mysql.Biz{
		Nickname: req.Nickname, Password: req.Pwd, LastLoginIp: c.ClientIP(),
	}); err != nil {
		model.Logger.Info(
			err.Error(),
			zap.Int("code", api.MsgRegisterCreateUserErr),
			zap.String("err", err.Error()),
		)
		api.JSON(c, api.MsgRegisterCreateUserErr, "create user is error")
		return
	}
	api.JSONOK(c)
	return
}

func Register(c *gin.Context) {
	req := &trans.ReqRegister{}
	if err := c.Bind(&req); err != nil {
		api.JSON(c, api.MsgErr, "request add member params is error")
		return
	}

	// 检查邮箱是否被注册
	info, e := service.Biz.InfoX(c, mysql.Conds{
		"email": req.Email,
	})
	if info.OpenId != 0 {
		api.JSON(c, api.MsgErr, "账号已经存在！")
		return
	}

	// 生成随机昵称
	nickname, _ := service.Biz.GetRandNickname()

	// 插入记录到数据库
	pwd, _ := util.Hash(req.Pwd)
	ret := mysql.Biz{0, "", nickname, "",
		"", "", time.Now().Unix(), time.Now().Unix(),
		pwd, req.Email, "", "", 1}
	e = service.Biz.Create(c, model.Db, &ret)
	if e != nil {
		api.JSON(c, api.MsgErr, e.Error())
		return
	}
	id := ret.OpenId
	model.Logger.Info("register info", zap.Any("id", id))

	// 生成并存储token
	token := funk.RandomString(16)
	cache.Cli.Set("fs:reg:"+cast.ToString(id), token, 300)

	// TODO 模板可能不需要传id
	data := map[string]interface{}{
		"user":   mysql.Biz{OpenId: id, Nickname: nickname, Email: req.Email},
		"token":  token,
		"domain": bootstrap.Conf.Server.DomainApi,
	}
	if e = service.Mailer.Send(service.RegisterEmail, req.Email, data); e != nil {
		model.Logger.Info(e.Error(), zap.Int("code", api.MsgErr))
		api.JSON(c, api.MsgErr, "邮件发送失败，请稍后再试！")
		return
	}
	api.JSONOK(c)
	return
}

func Confirm(c *gin.Context) {
	id := cast.ToInt(c.Query("id"))
	token := c.Query("token")
	relToken, _ := cache.Cli.GetString("fs:reg:" + cast.ToString(id))
	if token != relToken {
		api.JSON(c, api.MsgErr, "Token不合法，请重新注册！")
	}
	// 更新账号状态
	if e := service.Biz.Update(c, model.Db, id, mysql.Ups{
		"status": 0,
	}); e != nil {
		model.Logger.Info(e.Error(), zap.Int("code", api.MsgErr))
		api.JSON(c, api.MsgErr, "注册失败，请稍后再试！")
	}
	c.Redirect(302, "http://"+bootstrap.Conf.Server.DomainWeb+"/#/login")
}
