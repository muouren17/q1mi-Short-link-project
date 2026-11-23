# Shortener

基于 Go、go-zero 和 MySQL 的短链接服务，目前处于开发阶段。`POST /convert` 已实现链接可达性检查和去重查询，短码生成与保存仍待实现。`GET /:shortURL` 已注册路由，短链查询逻辑待实现。

## 本地启动

需要 Go 1.26 和 MySQL。新建数据库 `shortener`（也可以选用已有的开发数据库），然后导入 [`sequence.sql`](sequence.sql) 和 [`short_url_map.sql`](short_url_map.sql)。

```sh
mysql_host=127.0.0.1
mysql_user=your_mysql_user
mysql -h "$mysql_host" -P 3306 -u "$mysql_user" -p shortener < sequence.sql
mysql -h "$mysql_host" -P 3306 -u "$mysql_user" -p shortener < short_url_map.sql
```

复制配置模板，并填写 `SequenceURLDB.DSN` 和 `ShortURLDB.DSN` 中的账号、密码、地址和数据库名。两项可以使用同一个数据库。`etc/shortener-api.yaml` 由 Git 忽略，真实凭据填在该文件里。

```sh
if [ ! -e etc/shortener-api.yaml ]; then
  install -m 600 etc/shortener-api.example.yaml etc/shortener-api.yaml
fi
```

如果在 WSL 2 的 NAT 模式下运行服务，而 MySQL 在 Windows 上，把上述 `mysql_host` 和两个 DSN 的地址改为下面命令返回的 Windows 网关 IP。

```sh
ip -4 route show default | awk 'NR == 1 {print $3}'
```

配置完成后启动服务，模板中的监听地址为 `http://127.0.0.1:8888`。

```sh
go run . -f etc/shortener-api.yaml
```

`POST /convert` 会由服务端访问提交的 URL，错误日志可能包含完整的 URL。开发时请使用不含访问令牌或其他敏感信息的测试链接。

## 开发

- [`shortener.api`](shortener.api) 定义接口，`internal/logic/` 存放业务逻辑。
- `model/` 存放 MySQL 模型。安装 `goctl` 后，可用下列命令从 SQL 文件重新生成模型。

```sh
goctl model mysql ddl --src='./*.sql' --dir=./model
go test ./...
```
