// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package svc

import (
	"shortener/api/internal/config"
	"shortener/api/internal/repository"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	SequenceURLModel model.SequenceModel
	ShortURLModel    model.ShortUrlMapModel

	Sequence          sequence.Sequence
	ShortBlackListMap map[string]struct{}
	ShortURLRedis     repository.CachedShortURLModel
	Filter            *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqconn := sqlx.NewMysql(c.SequenceURLDB.DSN)
	shconn := sqlx.NewMysql(c.ShortURLDB.DSN)
	m := make(map[string]struct{}, len(c.ShortBlackList))
	for _, b := range c.ShortBlackList {
		m[b] = struct{}{}
	}
	bloomstore := redis.MustNewRedis(redis.RedisConf{
		Host:     c.Redis.Addr,
		Type:     redis.NodeType,
		Pass:     c.Redis.Password,
		NonBlock: true,
	})
	filter := bloom.New(bloomstore, "shortener:bloom:v1", 2_000_000)
	return &ServiceContext{
		Config:           c,
		SequenceURLModel: model.NewSequenceModel(sqconn),
		ShortURLModel:    model.NewShortUrlMapModel(shconn),
		// Sequence:         sequence.NewMySQL(c.SequenceURLDB.DSN),
		Sequence:          sequence.NewRedis(c.Redis.Addr, c.Redis.Password),
		ShortURLRedis:     *repository.NewCachedShortURLModel(c.Redis.Addr, c.Redis.Password),
		Filter:            filter,
		ShortBlackListMap: m,
	}
}
