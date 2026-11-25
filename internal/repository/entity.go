package repository

import "errors"

const (
	S_L_urlKey = "lwithsURL"
)

var ErrRedisNotFound = errors.New("缓存未命中")
