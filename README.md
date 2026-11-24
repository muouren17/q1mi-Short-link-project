# Shortener

Shortener 是一个基于 Go 和 go-zero 的短链接服务。服务目前已接入链接校验、MySQL 映射查询和 Redis 取号。

## HTTP 接口

接口定义见 [`shortener.api`](shortener.api)。

| 接口 | 输入 | 当前行为 |
| --- | --- | --- |
| `POST /convert` | JSON 字段 `longURL` | 校验链接并分配序号。 |
| `GET /:shortURL` | 路径参数 `shortURL` | 路由已注册，短链查询待实现。 |

`POST /convert` 先检查原链接能否返回 HTTP `200`，再查询 MySQL，按长链接的 MD5 判断是否已转换，并用 URL 路径的最后一段匹配已有短码。检查通过后调用 `Sequence.Next()`。短码生成、映射写入和短链接返回待实现。[转换逻辑](internal/logic/convertLogic.go)

## 序号分配

`sequence.Sequence` 通过 `Next() (uint64, error)` 提供统一的取号接口。[接口定义](sequence/sequence.go)

- **Redis（默认）** 使用单次原子操作 `INCR IncryCount`，直接以返回值作为序号。键缺失时首次取号为 `1`，应用重启后会沿用 Redis 中已有的计数。Redis 错误和负数结果都会作为取号错误返回。[Redis 实现](sequence/redis.go)
- **MySQL** 通过 `REPLACE INTO sequence` 写入固定的 `stub` 值，并读取自增 ID。切换时，将 `sequence.sql` 导入 `SequenceURLDB.DSN` 对应的数据库，并将服务上下文中的取号实现改为 `sequence.NewMySQL(c.SequenceURLDB.DSN)`。[表结构](sequence.sql) · [MySQL 实现](sequence/mysql.go) · [服务上下文](internal/svc/serviceContext.go)

## 本地启动

需要 Go 1.26、MySQL 和 Redis。默认配置只需创建 `shortener` 数据库并导入长短链映射表。使用已有开发库时，从导入命令开始，并将 `shortener` 改为实际库名。[映射表结构](short_url_map.sql)

```sh
mysql_host=127.0.0.1
mysql_user=your_mysql_user
mysql -h "$mysql_host" -P 3306 -u "$mysql_user" -p \
  -e 'CREATE DATABASE IF NOT EXISTS shortener CHARACTER SET utf8mb4;'
mysql -h "$mysql_host" -P 3306 -u "$mysql_user" -p shortener < short_url_map.sql
```

复制配置模板。`etc/shortener-api.yaml` 由 Git 忽略，真实凭据填在该文件里。[配置模板](etc/shortener-api.example.yaml)

```sh
if [ ! -e etc/shortener-api.yaml ]; then
  install -m 600 etc/shortener-api.example.yaml etc/shortener-api.yaml
fi
```

将 `ShortURLDB.DSN` 指向已导入映射表的数据库，并填写 `SequenceURLDB.DSN`。本地开发时两项可以使用同一个 DSN。在配置文件顶层加入以下 Redis 配置，按实际部署填写地址和密码。[配置字段](internal/config/config.go)

```yaml
Redis:
  Addr: "127.0.0.1:6379"
  Password: ""
```

如果服务运行在 WSL 2 的 NAT 模式下，而 MySQL 或 Redis 运行在 Windows 上，将相应的数据库地址换成下面命令返回的 Windows 网关 IP。

```sh
ip -4 route show default | awk 'NR == 1 {print $3}'
```

启动 MySQL 和 Redis 后运行服务。模板中的监听地址为 `http://127.0.0.1:8888`。

```sh
go run . -f etc/shortener-api.yaml
```

提交一个服务端可访问的测试链接。取号成功时，响应为 `200 OK`，正文为 JSON `null`，序号显示在服务终端。

```sh
curl -i http://127.0.0.1:8888/convert \
  -H 'Content-Type: application/json' \
  --data '{"longURL":"https://example.com/"}'
```

## 开发

- 接口业务逻辑位于 [`internal/logic/`](internal/logic/)。
- 取号接口及 Redis、MySQL 实现位于 [`sequence/`](sequence/)。
- MySQL 数据模型位于 [`model/`](model/)。
- URL 处理、HTTP 连通性检查和 MD5 摘要工具位于 [`pkg/`](pkg/)。

安装 `goctl` 后，可用下列命令从 SQL 文件重新生成模型。

```sh
goctl model mysql ddl --src='./*.sql' --dir=./model
```

## 测试

```sh
go test ./...
```

`pkg/connect` 的测试会尝试访问外部网站。只运行 URL 解析和 MD5 测试时，可执行 `go test ./pkg/urltool ./pkg/md5`。
