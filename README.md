# MetarService

[![ReleaseCard]][Release]![ReleaseDataCard]![LastCommitCard]  
![BuildStateCard]![DockerStateCard]![ProjectLicense]

MetarService 是一个METAR数据服务，是微服务架构的一部分，用于获取METAR数据

本项目可以独立运行，单独提供METAR数据服务

本项目提供了RESTful API接口与GRPC接口，用于获取 METAR 数据

## Feature列表

- [X] 格式化获取METAR数据
- [ ] 解析METAR数据
- [X] 格式化获取TAF数据
- [ ] 解析TAF数据

## 如何使用

### ***(推荐)*** 使用Docker部署

1. ***(推荐)*** 使用docker-compose部署  
   i. 克隆或下载本项目到本地，并进入`docker`目录  
   ii. 按需编辑配置文件或`docker-compose.yml`文件  
   iii. 运行`docker-compose up -d`命令  
   iv. 访问[http://127.0.0.1:8080](http://127.0.0.1:8080)查看是否部署成功  
   v. 如果需要添加命令行参数
   ```yml
   services:
     fsd:
       image: halfnothing/metar-service:latest
       # 省略部分字段
       command:
         - "-thread 32"
   ```
   推荐使用环境变量代替命令行参数
   ```yml
   services:
     fsd:
       image: halfnothing/metar-service:latest
       # 省略部分字段
       environment:
         QUERY_THREAD: 32
   ```

2. 使用docker命令部署  
   命令示例如下
   ```shell
   docker run -d --name metar-service -p 8080:8080 -v $(pwd)/config.yaml:/metar-service/config.yaml halfnothing/metar-service:latest
   ``` 
   如果需要添加命令行参数, 则在命令的最后添加
   ```shell
   docker run -d ... halfnothing/metar-service:latest -thread 32
   ```

3. 通过Dockerfile构建  
   i. 手动构建
   ```shell
   # 克隆本仓库
   git clone https://github.com/FSD-Universe/metar-service.git
   # 进入项目目录
   cd metar-service
   # 运行docker构建
   docker build -t metar-service:latest .
   # 运行docker容器
   docker run -d --name metar-service -p 8080:8080 -v $(pwd)/config.yaml:/metar-service/config.yaml metar-service:latest
   ```
   ii. 自动构建
   ```shell
   # 克隆本仓库
   git clone https://github.com/FSD-Universe/metar-service.git
   # 进入项目目录
   cd metar-service
   # 进入docker目录并且修改docker-compose.yml文件
   cd docker
   vi docker-compose.yml
   ```
   将`image: halfnothing/metar-service:latest`这一行替换为`build: ".."`    
   然后在同目录运行
   ```shell
   docker compose up -d
   ```

### 普通部署

1. 获取项目可执行文件
    - 前往[Release]页面下载最新版本
    - 前往[Action]页面下载最新开发版本
    - 手动[编译](#手动构建)本项目
2. [可选]下载[`config.yaml`](./docker/config.yaml)配置文件放置于可执行文件同级目录中
3. 运行可执行文件，如果配置文件存在，则使用配置文件，否则创建默认配置文件

## 手动构建

```shell
# 克隆本仓库
git clone https://github.com/FSD-Universe/metar-service.git
# 进入项目目录
cd metar-service
# 确认安装了go编译器并且版本>=1.24.6
go version
# 运行go build命令
go build -ldflags="-w -s" -tags "http telemetry" .
# 对于windows系统, 可执行文件为metar-service.exe
# 对于linux系统, 可执行文件为metar-service
# [可选]使用upx压缩可执行文件
# windows
upx.exe -9 metar-service.exe
# linux
upx -9 metar-service
```

## 命令行参数与环境变量一览

| 命令行参数                 | 环境变量                  | 描述                 | 默认值                                       |
|:----------------------|:----------------------|:-------------------|:------------------------------------------|
| no_logs               | NO_LOGS               | 禁用日志输出到文件          | false                                     |
| auto_migrate          | AUTO_MIGRATE          | 自动迁移数据库(不要在生产环境使用) | false                                     |
| config                | CONFIG_FILE_PATH      | 配置文件路径             | "config.yaml"                             |
| health_check_interval | HEALTH_CHECK_INTERVAL | 健康检查间隔             | "30s"                                     |
| health_check_timeout  | HEALTH_CHECK_TIMEOUT  | 健康检查超时时间           | "5s"                                      |
| deregister_after      | DEREGISTER_AFTER      | 健康检查失败后注销时间        | "1m"                                      |
| service_address       | SERVICE_ADDRESS       | 服务对外访问地址，默认为本地网卡地址 | "localhost"                               |
| center_address        | CENTER_ADDRESS        | consul注册中心地址       | "localhost:8500"                          |
| reconnect_timeout     | RECONNECT_TIMEOUT     | 重连超时时间             | "30s"                                     |
| eth_name              | ETH_NAME              | 以太网接口名称            | "Ethernet"(windows) / "eth0"(linux/macos) |
| http_timeout          | HTTP_TIMEOUT          | Http请求超时时间         | "30s"                                     |
| gzip_level            | GZIP_LEVEL            | Gzip压缩等级           | 5                                         |

## 贡献指南

1. 开一个 Issue 与我们讨论
2. Fork 本项目并完成你的修改
3. 不要修改任何除了你创建以外的源代码的版权信息
4. 遵守良好的代码编码规范
5. 开一个 Pull Request

## 开源协议

MIT License

Copyright © 2025 Half_nothing

无附加条款。

[ReleaseCard]: https://img.shields.io/github/v/release/FSD-Universe/metar-service?logo=github&style=for-the-badge

[ReleaseDataCard]: https://img.shields.io/github/release-date/FSD-Universe/metar-service?display_date=published_at&logo=github&style=for-the-badge

[LastCommitCard]: https://img.shields.io/github/last-commit/FSD-Universe/metar-service?display_timestamp=committer&logo=github&style=for-the-badge

[BuildStateCard]: https://img.shields.io/github/actions/workflow/status/FSD-Universe/metar-service/go-build.yml?logo=go&label=Build&style=for-the-badge

[DockerStateCard]: https://img.shields.io/github/actions/workflow/status/FSD-Universe/metar-service/push-latest.yml?logo=docker&label=Push&style=for-the-badge

[ProjectLanguageCard]: https://img.shields.io/github/languages/top/FSD-Universe/metar-service?logo=github&style=for-the-badge

[ProjectLicense]: https://img.shields.io/badge/License-MIT-blue?logo=github&style=for-the-badge

[Release]: https://www.github.com/FSD-Universe/metar-service/releases/latest

[Action]: https://github.com/FSD-Universe/metar-service/actions/workflows/go-build.yml

[Release]: https://www.github.com/FSD-Universe/metar-service/releases/latest

