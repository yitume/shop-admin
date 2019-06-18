package service

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"git.yitum.com/saas/shop-admin/model"
	"git.yitum.com/saas/shop-admin/model/mysql"
)

type app struct{}

func InitApp() *app {
	return &app{}
}

func (*app) List(currentPage, pageSize int, query string, sort string) (total int, appList []*mysql.App) {
	sql := model.Db.Table("app")
	appList = make([]*mysql.App, 0)
	if aid, err := strconv.Atoi(query); err == nil && aid != 0 {
		sql = sql.Where("aid = ?", aid)
	} else if s := strings.TrimSpace(query); s != "" {
		sql = sql.Where("name LIKE ?", "%"+s+"%")
	}
	sql.Count(&total)
	offset := 0
	offset = currentPage * pageSize
	sql.Order(sort).Offset(offset).Limit(pageSize).Find(&appList)
	return
}

func (*app) Add(name, redirectUri string) (err error) {
	create := mysql.App{
		Name:        name,
		Secret:      getRandomString(32),
		RedirectUri: redirectUri,
		CallNo:      0,
		Status:      1,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	if err = model.Db.Create(&create).Error; err != nil {
		model.Logger.Debug("add user create error", zap.String("err", err.Error()))
		return
	}
	return nil
}

func (*app) Update(aid int, name, redirectUri string, status int) (err error) {
	if err = model.Db.Table("app").Where("aid=?", aid).Updates(mysql.Ups{
		"name":         name,
		"redirect_uri": redirectUri,
		"status":       status,
	}).Error; err != nil {
		model.Logger.Debug("add user update error", zap.String("err", err.Error()))
		return
	}
	return nil
}

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
