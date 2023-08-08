<!-- # react-dashboard -->
## React 后台管理系统

基于 `React` 生态系统搭建的后台管理系统模板

如果觉得不错或者对你有些许的帮助，欢迎 `star`，或者你有更好的实现方式、有趣的 `idea`，也欢迎留言交流 

以下为项目基本信息↓

### 技术栈

`React@18.2.0 + React-Router@6.14.2 + Antd@5.8.2`

> `Create React App`    脚手架工具快速搭建项目结构

> `braft-editor@2.3.8`    富文本插件

> `echarts@4.4.0`   数据可视化

### 基本功能

- [x] 登陆/注销功能
- [x] axios 封装
- [x] 简单权限管理

### 项目结构

```
├── public                   # 不参与编译的资源文件
├── src                      # 主程序目录
│   ├── request                     # axios 封装
│   ├── assets                  # 资源文件
│   │   └── images                  # 图片资源
│   ├── components              # 全局公共组件
│   │   ├── AuthGuard        # 路由鉴权
│   │   ├── Redirect        # 路由重定向
│   │   └── ProtectedLayout              # 登录后侧边栏菜单
│   ├── hooks             # 自定义钩子
│   │   ├── useAuth        # 鉴权、管理员信息、登入、登出
│   │   └── useLocalStorage         # 本地存储
│   ├── styles                   # 样式目录
│   ├── utils                   # 工具类
│   ├── views                   # UI 页面
│   ├── APP.js                  # App.js
│   └── index.js                # index.js
```

### 使用方法

```npm
git clone https://github.com/Oaciqihz/react-dashboard.git

cd react-dashboard

// 安装依赖
yarn

// 启动
yarn start

// 打包
yarn build

```