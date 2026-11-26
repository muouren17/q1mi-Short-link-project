// Code scaffolded by goctl. Safe to edit.
// goctl 1.11.0

package main

import (
	"flag"
	"fmt"

	"shortener/api/internal/config"
	"shortener/api/internal/handler"
	"shortener/api/internal/svc"
	"shortener/pkg/base62"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shortener-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	base62.MustInit(c.Base62Str)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
