package service

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/mysql"
)

func (*area) GetAddr(c *gin.Context, aid int) (areas []mysql.Area, err error) {
	area, err := Area.InfoX(c, mysql.Conds{
		"id": aid,
	})
	if err != nil {
		return
	}
	city, err := Area.InfoX(c, mysql.Conds{
		"id": area.Pid,
	})
	if err != nil {
		return
	}
	province, err := Area.InfoX(c, mysql.Conds{
		"id": city.Pid,
	})
	if err != nil {
		return
	}
	return []mysql.Area{area, city, province}, nil
}
