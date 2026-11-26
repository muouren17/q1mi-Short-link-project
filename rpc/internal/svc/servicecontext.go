package svc

import (
	"shortener/model"
	"shortener/rpc/internal/config"
	"shortener/rpc/internal/repository"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config            config.Config
	SequenceURLModel  model.SequenceModel
	ShortURLModel     model.ShortUrlMapModel
	ShortURLRedis     repository.CachedShortURLModel
	Sequence          sequence.Sequence
	ShortBlackListMap map[string]struct{}
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
		Host:     c.RRedis.Addr,
		Type:     redis.NodeType,
		Pass:     c.RRedis.Password,
		NonBlock: true,
	})
	filter := bloom.New(bloomstore, "shortener:bloom:v1", 2_000_000)
	return &ServiceContext{
		Config:           c,
		SequenceURLModel: model.NewSequenceModel(sqconn),
		ShortURLModel:    model.NewShortUrlMapModel(shconn),
		// Sequence:         sequence.NewMySQL(c.SequenceURLDB.DSN),
		ShortURLRedis:     *repository.NewCachedShortURLModel(c.RRedis.Addr, c.RRedis.Password),
		Sequence:          sequence.NewRedis(c.RRedis.Addr, c.RRedis.Password),
		Filter:            filter,
		ShortBlackListMap: m,
	}
}
