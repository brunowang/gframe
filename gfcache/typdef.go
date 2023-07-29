package gfcache

import "errors"

var (
	RecordNotFound = errors.New("record not found")
	CacheNotFound  = errors.New("cache not found")
)
