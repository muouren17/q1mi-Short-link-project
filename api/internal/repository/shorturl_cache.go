package repository

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type CachedShortURLModel struct {
	client *redis.Client
	g      *singleflight.Group
}

func NewCachedShortURLModel(addr string, password string) *CachedShortURLModel {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &CachedShortURLModel{
		client: c,
		g:      new(singleflight.Group),
	}
}

func (c *CachedShortURLModel) FindOneBySurl(ctx context.Context, surl string) (lurl string, err error) {
	if len(surl) == 0 {
		return "", nil
	}
	lurl, err = c.client.HGet(S_L_urlKey, surl).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrRedisNotFound
		}
		return "", err
	}
	return lurl, nil
}

func (c *CachedShortURLModel) SaveUrl(ctx context.Context, surl string, lurl string) (bool, error) {
	if len(lurl) == 0 || len(surl) == 0 {
		return false, nil
	}
	ok, err := c.client.HSet(S_L_urlKey, surl, lurl).Result()
	return ok, err
}

func (c *CachedShortURLModel) GetOrLoad(ctx context.Context, surl string, load func(context.Context) (string, error)) (string, error) {
	if lurl, err := c.FindOneBySurl(ctx, surl); err == nil && len(lurl) != 0 {
		return lurl, nil
	}
	ch := c.g.DoChan(surl, func() (any, error) {
		sharedCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		lurl, err := c.FindOneBySurl(sharedCtx, surl)
		if err == nil && len(lurl) != 0 {
			return lurl, nil
		}
		if err != nil && !errors.Is(err, ErrRedisNotFound) {
			logx.Errorw("ShortURLRedis.FindOneBySurl failed", logx.LogField{Key: "err", Value: err.Error()})
		}
		lurl, err = load(sharedCtx)
		if err != nil {
			return nil, err
		}
		_, err = c.SaveUrl(sharedCtx, surl, lurl)
		if err != nil {
			logx.Errorw("ShortURLRedis.SaveUrl failed", logx.LogField{Key: "err", Value: err.Error()})
		}
		return lurl, nil

	})
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case ret := <-ch:
		if ret.Err != nil {
			return "", ret.Err
		}
		return ret.Val.(string), nil
	}
}
