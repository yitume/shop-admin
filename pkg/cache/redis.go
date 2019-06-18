package cache

import (
	"github.com/yitume/caller/redigo"
)

var (
	Cli *redigo.RedigoClient
)

func InitRedis() {
	Cli = redigo.Caller("default")
}
