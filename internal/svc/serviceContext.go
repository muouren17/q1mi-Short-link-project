// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package svc

import (
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config           config.Config
	SequenceURLModel model.SequenceModel
	ShortURLModel    model.ShortUrlMapModel

	Sequence sequence.Sequence
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqconn := sqlx.NewMysql(c.SequenceURLDB.DSN)
	shconn := sqlx.NewMysql(c.ShortURLDB.DSN)
	return &ServiceContext{
		Config:           c,
		SequenceURLModel: model.NewSequenceModel(sqconn),
		ShortURLModel:    model.NewShortUrlMapModel(shconn),
		// Sequence:         sequence.NewMySQL(c.SequenceURLDB.DSN),
		Sequence: sequence.NewRedis(c.Redis.Addr, c.Redis.Password),
	}
}
