package mdw

import (
	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/router/mdw/auth"
)

func OpenId(c *gin.Context) int {
	return auth.Default(c).Id
}
