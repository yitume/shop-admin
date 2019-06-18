package service

import (
	"fmt"
	"time"

	"git.yitum.com/saas/shop-admin/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
)

const (
	UserInit    = 0
	UserActived = 1
	UserBanned  = 2
)

func (g *biz) GetRandNickname() (n string, err error) {
	n = "fs_" + funk.RandomString(8)
	var req = mysql.Biz{}
	if err = model.Db.Table("biz").Where("`nickname`=? ", n).First(&req).Error; err != nil {
		model.Logger.Error("biz update error", zap.String("err", err.Error()))
		return
	}
	if req.OpenId != 0 {
		return g.GetRandNickname()
	}
	return
}

func (*biz) AddBiz(nickname string, pwd string, ip string) (err error) {
	var pwdHash string
	pwdHash, err = util.Hash(pwd)
	if err != nil {
		model.Logger.Debug("add user hash error", zap.String("err", err.Error()))
		return
	}
	user := mysql.Biz{
		Nickname:    nickname,
		Password:    pwdHash,
		Status:      UserActived,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
		LastLoginIp: ip,
	}
	if err = model.Db.Create(&user).Error; err != nil {
		model.Logger.Debug("add user create error", zap.String("err", err.Error()))
		return
	}
	return nil
}

func (*biz) GetBizByPwd(nickname string, email string, pwd string, clientIp string) (resp *mysql.Biz, err error) {
	fmt.Println(nickname, email, pwd, clientIp)
	query := "1=1"
	data := make([]interface{}, 0)
	if nickname != "" {
		query += " and nickname = ? "
		data = append(data, nickname)
	}
	if email != "" {
		query += " and email = ? "
		data = append(data, email)
	}
	resp = &mysql.Biz{}
	if err = model.Db.Where(query, data).First(resp).Error; err != nil {
		model.Logger.Debug("GetBizByPwd ERROR", zap.String("err", err.Error()))
		return
	}
	err = util.Verify(resp.Password, pwd)
	if err != nil {
		model.Logger.Debug("verify error1", zap.String("err", err.Error()))
		return
	}
	if err = model.Db.Table("biz").Where("open_id = ?", resp.OpenId).Updates(gin.H{
		"updated_at":    time.Now().Unix(),
		"last_login_ip": clientIp,
	}).Error; err != nil {
		model.Logger.Debug("update user create error", zap.String("err", err.Error()))
		return
	}
	return
}

func (*biz) GetBizByOidPwd(openId int, pwd string) (resp *mysql.Biz, err error) {
	resp = &mysql.Biz{}
	if err = model.Db.Where("open_id = ?", openId).First(resp).Error; err != nil {
		model.Logger.Debug("GetBizByPwd ERROR", zap.String("err", err.Error()))
		return
	}
	err = util.Verify(resp.Password, pwd)
	if err != nil {
		model.Logger.Debug("verify error1", zap.String("err", err.Error()))
		return
	}
	return
}

func (*biz) UpdatePwd(openId int, pwd string) (err error) {
	var pwdHash string
	pwdHash, err = util.Hash(pwd)
	if err != nil {
		model.Logger.Debug("update user hash error", zap.String("err", err.Error()))
		return
	}
	if err = model.Db.Table("biz").Where("open_id = ?", openId).Updates(gin.H{
		"password":   pwdHash,
		"updated_at": time.Now().Unix(),
	}).Error; err != nil {
		model.Logger.Debug("update user create error", zap.String("err", err.Error()))
		return
	}
	return nil
}
