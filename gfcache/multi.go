package gfcache

type Cache interface {
	HasCache(key string) bool
	GetCache(key string, data interface{}) error
	SetCache(key string, data interface{}) error
	DelCache(key string) error
}

type MultiCache struct {
	caches []Cache
}

func NewMultiCache(caches ...Cache) *MultiCache {
	return &MultiCache{caches: caches}
}

func (c *MultiCache) Append(cache Cache) *MultiCache {
	c.caches = append(c.caches, cache)
	return c
}

func (c *MultiCache) HasCache(key string) bool {
	length := len(c.caches)
	if length == 0 {
		return false
	}
	return c.caches[length-1].HasCache(key)
}

func (c *MultiCache) GetCache(key string, data interface{}) error {
	length := len(c.caches)
	if length == 0 {
		return CacheNotFound
	}
	if !c.caches[length-1].HasCache(key) {
		_ = c.DelCacheFromLevel(key, length-2)
		return CacheNotFound
	}
	for i, cache := range c.caches {
		if err := cache.GetCache(key, data); err == nil {
			if i > 0 {
				_ = c.caches[i-1].SetCache(key, data)
			}
			return nil
		} else if err == RecordNotFound {
			if i > 0 {
				_ = c.caches[i-1].SetCache(key, nil)
			}
			return RecordNotFound
		}
	}
	return CacheNotFound
}

func (c *MultiCache) SetCache(key string, data interface{}) error {
	length := len(c.caches)
	if length == 0 {
		return nil
	}
	return c.caches[length-1].SetCache(key, data)
}

func (c *MultiCache) DelCache(key string) error {
	length := len(c.caches)
	if length == 0 {
		return nil
	}
	return c.caches[length-1].DelCache(key)
}

func (c *MultiCache) DelCacheFromLevel(key string, levelIdx int) error {
	if levelIdx < 0 {
		return nil
	}
	length := len(c.caches)
	if levelIdx > length-1 {
		levelIdx = length - 1
	}
	for i := levelIdx; i >= 0; i-- {
		err := c.caches[i].DelCache(key)
		if err != nil {
			return err
		}
	}
	return nil
}
