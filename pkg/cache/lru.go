package cache

import (
	"github.com/coocood/freecache"
)

var (
	LRU *freecache.Cache
)

func InitLRU() {
	LRU = freecache.NewCache(100 * 1024 * 1024)
}
