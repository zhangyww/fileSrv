/**
 * @Author: zhangyw
 * @Description:
 * @File:  CacheRepo
 * @Date: 2021/5/21 16:17
 */

package cacheRepo

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

var repoInstance = newCacheRepo(100000)

type CacheRepo struct {
	cache *ristretto.Cache
}

func newCacheRepo(maxCount int64) *CacheRepo {
	numCounter := maxCount * 10
	bufferItems := int64(64)
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounter,  // number of keys to track frequency of.
		MaxCost:     maxCount,    // maximum cost of cache.
		BufferItems: bufferItems, // number of keys per Get buffer.
	})
	return &CacheRepo{cache: cache}
}

func (this *CacheRepo) Get(key interface{}) (interface{}, bool) {
	return this.cache.Get(key)
}

func (this *CacheRepo) Set(key interface{}, value interface{}) {
	this.cache.Set(key, value, 1)
}

func (this *CacheRepo) SetWithTTL(key interface{}, value interface{}, ttl time.Duration) {
	this.cache.SetWithTTL(key, value, 1, ttl)
}

func (this *CacheRepo) Del(key interface{}) bool {
	this.cache.Del(key)
	return true
}

func GetRepo() *CacheRepo {
	return repoInstance
}
