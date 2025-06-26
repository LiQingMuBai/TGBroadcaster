package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	DefaultExpiration = 50 * time.Minute // 默认缓存过期时间30分钟
	CleanupInterval   = 1 * time.Hour    // 清理过期缓存项的间隔
)

var (
	// 全局缓存实例
	Cache = cache.New(DefaultExpiration, CleanupInterval)
)

// Set 设置缓存，使用默认过期时间
func Set(key string, value interface{}) {
	Cache.Set(key, value, cache.DefaultExpiration)
}

// SetWithExpiration 设置缓存并指定过期时间
func SetWithExpiration(key string, value interface{}, expiration time.Duration) {
	Cache.Set(key, value, expiration)
}

// Get 获取缓存值
func Get(key string) (interface{}, bool) {
	return Cache.Get(key)
}

// Delete 删除缓存项
func Delete(key string) {
	Cache.Delete(key)
}

// Flush 清空所有缓存
func Flush() {
	Cache.Flush()
}
