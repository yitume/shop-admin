package service

import (
	"git.yitum.com/saas/shop-admin/pkg/cache"
)

var (
	Storage *storage
	App     *app
	Mailer  *mailer
)

func Init() {
	cache.Init()
	Storage = NewStorage()
	App = InitApp()
	Mailer = InitMailer()
}
