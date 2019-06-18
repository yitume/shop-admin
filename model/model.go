package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/yitume/caller/ginsession"
	cgorm "github.com/yitume/caller/gorm"
	"github.com/yitume/caller/zap"
)

var (
	Db      *gorm.DB
	Logger  *zap.ZapClient
	Session gin.HandlerFunc
)

func Init() {
	Db = cgorm.Caller("oauth").DB
	Logger = zap.Caller("system")
	Session = ginsession.Caller()
}
