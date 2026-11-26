package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	SequenceURLDB struct {
		DSN string
	}
	ShortURLDB struct {
		DSN string
	}
	RRedis struct {
		Addr     string
		Password string
	}
	Base62Str      string
	ShortBlackList []string
	ShortDomain    string
}
