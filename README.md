# Shortener

Shortener 是一个基于 Go 和 go-zero 的短链接服务，支持创建短链和访问短链跳转。MySQL 保存长短链映射，Redis 用于默认取号和缓存查询结果。

## HTTP 接口

接口定义见 [`shortener.api`](shortener.api)。

| 接口 | 输入 | 当前行为 |
| --- | --- | --- |
| `POST /convert` | JSON 字段 `longURL` | 校验长链接，生成并保存短码，返回 `shortURL`。 |
| `GET /:shortURL` | 路径中的短码 | 查询原链接，返回 `302 Found`，并在 `Location` 响应头中给出跳转地址。 |

创建时，服务先检查长链接能否返回 HTTP `200`，再用长链接的 MD5 查询是否已有映射。通过检查后，服务从 Redis 取号，用 Base62 编码成短码，将映射写入 MySQL，并返回由 `ShortDomain` 和短码组成的地址。黑名单配置字段已加入，过滤尚待接入服务上下文。[转换逻辑](internal/logic/convertLogic.go) · [服务上下文](internal/svc/serviceContext.go)

## 跳转与缓存

访问短链时，服务先从 Redis 哈希 `lwithsURL` 按短码查找原链接。缓存未命中时，同一服务进程内并发访问同一短码的请求由 `singleflight` 合并。共享查询会再检查一次缓存，仍未命中时查询 MySQL 并回填 Redis。[跳转逻辑](internal/logic/showLogic.go) · [缓存实现](internal/repository/shorturl_cache.go)

Redis 读取出错时，服务会尝试从 MySQL 查询原链接。只要查到了原链接，即使缓存回填失败，服务也会返回 `302 Found`。

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

将 `ShortURLDB.DSN` 指向已导入映射表的数据库，并填写 `SequenceURLDB.DSN`。本地开发时两项可以使用同一个 DSN。在配置文件顶层加入以下字段，并按实际环境填写 Redis 地址和密码。[配置定义](internal/config/config.go)

```yaml
Redis:
  Addr: "127.0.0.1:6379"
  Password: ""

Base62Str: "qJm5WH9zfQ2aCvRkZOo0Vjx8Lp3PsTK7ndEB1utXcS6FhyYwM4ArNgeUIlDbiG"
ShortBlackList: ["api", "convert"]
ShortDomain: "http://127.0.0.1:8888"
```

`Base62Str` 使用由 `0-9`、`a-z` 和 `A-Z` 组成的 62 个互不重复的字符。`ShortDomain` 包含协议、主机和可选端口，末尾省略 `/`，服务会在拼接时补上。

如果服务运行在 WSL 2 的 NAT 模式下，而 MySQL 或 Redis 运行在 Windows 上，将相应的数据库地址换成下面命令返回的 Windows 网关 IP。

```sh
ip -4 route show default | awk 'NR == 1 {print $3}'
```

启动 MySQL 和 Redis 后运行服务。模板中的监听地址为 `http://127.0.0.1:8888`。

```sh
go run . -f etc/shortener-api.yaml
```

## 创建和访问短链

提交一个服务端可访问的长链接，验证创建接口。

```sh
curl -i http://127.0.0.1:8888/convert \
  -H 'Content-Type: application/json' \
  --data '{"longURL":"https://example.com/"}'
```

使用上面的字符表且 Redis 首次取号时，成功响应为 `200 OK`，正文如下。短码随当前序号变化。

```json
{"shortURL":"http://127.0.0.1:8888/J"}
```

将命令中的地址换成实际返回的 `shortURL`。响应会包含 `302 Found` 和 `Location: https://example.com/`。浏览器打开这个地址时会自动跳转。

```sh
curl -sS -D - -o /dev/null http://127.0.0.1:8888/J
```

## 开发

- 接口业务逻辑位于 [`internal/logic/`](internal/logic/)。
- Redis 缓存和并发查询合并位于 [`internal/repository/`](internal/repository/)。
- 取号接口及 Redis、MySQL 实现位于 [`sequence/`](sequence/)。
- MySQL 数据模型位于 [`model/`](model/)。
- Base62 编解码位于 [`pkg/base62/`](pkg/base62/)。
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
