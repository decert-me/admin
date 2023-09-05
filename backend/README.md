# backend
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