package utils

import (
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

func init() {
	c = cache.New(2*time.Minute, 3*time.Minute)
}

func GetCache(key string) (any, bool) {
	return c.Get(key)
}

func SetCache(key string, val any) {
	c.Set(key, val, cache.DefaultExpiration)
}
func DelCache(key string) {
	c.Delete(key)
}

func DeleteWithPrefix(prefix string) {
	// 遍历所有缓存项
	for key := range c.Items() {
		// 如果键以指定前缀开头，则删除该缓存项
		if strings.HasPrefix(key, prefix) {
			c.Delete(key)
		}
	}
}
