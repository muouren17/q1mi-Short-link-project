// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	SequenceURLDB struct {
		DSN string
	}
	ShortURLDB struct {
		DSN string
	}
	Redis struct {
		Addr     string
		Password string
	}
}
