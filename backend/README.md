# backend

Decert 项目后台接口

## 运行环境

```shell
- Golang >= v1.19
- PostgreSQL
```

## 安装
```bash
git clone https://github.com/decert-me/admin.git
cd admin/backend
```
## 编译
```bash
# 主程序
go build -o decert_admin
```
## 配置
```bash
cp ./config/config.demo.yaml ./config/config.yaml
vi ./config/config.yaml
```
## 运行
```bash
# 主程序
./decert_admin
```

## 配置说明

### 运行端口配置

配置项：

```yaml
# system configuration
system:
  env: develop
  addr: 9092
```

env：运行环境，可选值为 develop、test、production

addr：运行端口

### 数据库配置

配置项：
```yaml
# pgsql configuration
pgsql:
  path: "127.0.0.1"
  port: "5432"
  config: ""
  db-name: ""
  username: "postgres"
  password: "123456"
  auto-migrate: true
  prefix: ""
  slow-threshold: 200
  max-idle-conns: 10
  max-open-conns: 100
  log-mode: "info"
  log-zap: false
```
path：数据库地址

port：数据库端口

config：数据库配置

db-name：数据库名称

username：数据库用户名

password：数据库密码

auto-migrate：是否自动迁移数据库

prefix：数据库表前缀

slow-threshold：慢查询阈值，单位毫秒

max-idle-conns：最大空闲连接数

max-open-conns：最大连接数

log-mode：日志级别

log-zap：是否使用zap日志库

### 日志级别配置

配置项：
```yaml
# log configuration
log:
  level: info
  save: true
  format: console
  log-in-console: true
  prefix: '[backend-go]'
  director: log
  show-line: true
  encode-level: LowercaseColorLevelEncoder
  stacktrace-key: stacktrace
```

level：日志级别 debug、info、warn、error、dpanic、panic、fatal

save：是否保存日志

format：日志格式

log-in-console：是否在控制台输出日志

prefix：日志前缀

director：日志保存路径

show-line：是否显示行号

encode-level：日志编码级别

stacktrace-key：堆栈信息

### JWT 配置

配置项：

```yaml
# auth configuration
auth:
  signing-key: "Decert"
  expires-time: 86400
  issuer: "Decert"
```

signing-key：签名密钥

expires-time：过期时间，单位秒

issuer：签发人


### 文件上传配置

配置项：

```yaml
# local configuration
local:
  path: 'uploads/file'
```

path：上传文件保存路径


### Casbin 配置

配置项：

```yaml
# casbin configuration
casbin:
  model-path: "assets/rbac_model.conf"
```
model-path: 配置文件路径

### 挑战信息配置

配置项：

```yaml
# quest configuration
quest:
  encrypt-key: "eb5a5bb2-ebbd-45cc-9d37-77a9377f2aca"
```

encrypt-key：挑战信息加密密钥

### IPFS 配置

配置项：

```yaml
# ipfs configuration
ipfs:
  - api: "https://ipfs.io/ipfs"
    upload-api: "http://192.168.1.10:3022/v1"
```

api：IPFS节点地址

upload-api：IPFS上传API地址[ipfs-uploader](https://e.gitee.com/upchain/repos/upchain/ipfs-uploader/sources)

### 打包配置

打包教程文档

配置项：
```yaml
pack:
  server: "http://192.168.1.26:8889"                 # 打包模块URL
  publish-path: "/Users/mac/Code/tutorials/publish"   # 打包后的文件存放路径（发布目录）
```

server：打包模块[backend-pack](https://github.com/decert-me/backend-pack)接口URL

publish-path: 打包后的文件存放路径（发布目录）

### 空投配置

空投NFT

配置项：
```yaml
airdrop:
  verify-key: "123456"                   # 校验key
  api: "http://192.168.1.10:8105"        # 回调接口
```

verify-key：[airdrop-backend](https://e.gitee.com/upchain/repos/upchain/airdrop-backend/sources)项目配置的verify-key

api: [airdrop-backend](https://e.gitee.com/upchain/repos/upchain/airdrop-backend/sources)项目的接口URL

### 挑战国际化配置

挑战国际化处理，Github项目：[document](https://github.com/decert-me/document)

配置项：
```yaml
# translate configuration
translate:
  github-repo: "https://github.com/decert-me/document"
  github-branch: "main"
```

github-repo：Github 挑战翻译的 repo

github-branch: Github repo 所在分支

### Zcloak 证书生成配置

生成zcloak证书

配置项：
```yaml
# ZCloak configuration
zcloak:
  url: "http://192.168.1.10:4030"
```

url：zcloak证书生成服务API URL

### NFT 认证页面接口


配置项：
```yaml
# nft configuration
nft:
  api: "http://192.168.1.10:8093/v1"
  api-key: "test"
```

api: [nft-collect](https://github.com/decert-me/nft-collect) 项目的API URL

api-key: [nft-collect](https://github.com/decert-me/nft-collect) 项目配置的API KEY


