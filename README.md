# VPT - Private Tracker System

Go 语言实现的私有 BT Tracker 系统，由 4 个微服务 + 1 个共享库组成。

## 架构

```
                 ┌─────────────┐
                 │   Client    │ (HTTP / UDP)
                 └──────┬──────┘
                        │
                 ┌──────▼──────┐
                 │   Gateway   │  网关：转发 + 鉴权 (HTTP/UDP)
                 └──────┬──────┘
        ┌───────────────┼───────────────┐
   ┌────▼────┐    ┌─────▼─────┐   ┌─────▼─────┐
   │ Admin   │    │  Tracker  │   │ Registry  │
   │(用户/种子)   │(HTTP/UDP+统计)  │(注册+配置) │
   └────┬────┘    └─────┬─────┘   └───────────┘
     sqlite           sqlite          sqlite
```

## 服务列表

| 服务     | 端口  | 职责                                              |
| -------- | ----- | ------------------------------------------------- |
| registry | 8500  | 服务注册发现 + 配置中心                           |
| gateway  | 8000  | HTTP/UDP 网关，请求转发 + 鉴权                    |
| admin    | 8001  | 用户注册/登录/鉴权 + BT 种子上传管理              |
| tracker  | 8002  | BT Tracker 协议 (HTTP/UDP) + 数据统计             |

## 服务间通信

- gateway -> admin：调用 `POST /api/v1/auth/verify` 完成请求鉴权
- 所有服务启动时向 registry 注册，并定期拉取配置
- 服务间 RPC 统一使用 HTTP + JSON

## 目录结构

```
vpt/
├── go.work
├── common/         共享库：日志、HTTP 客户端、注册客户端、中间件
├── registry/       注册 + 配置中心
├── gateway/        网关
├── admin/          管理中心 + 用户中心
└── tracker/        Tracker 服务
```

每个服务内部分层：

```
service/
├── cmd/main.go     启动入口
├── internal/
│   ├── config/     配置加载
│   ├── handler/    HTTP/UDP handler
│   ├── service/    业务逻辑
│   ├── repo/       sqlite 数据访问
│   └── model/      领域模型
├── data/           sqlite 数据文件
└── go.mod
```

## 启动

```bash
cd registry && go run ./cmd
cd gateway  && go run ./cmd
cd admin    && go run ./cmd
cd tracker  && go run ./cmd
```
