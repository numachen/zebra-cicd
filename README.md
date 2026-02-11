# zebra-cicd


## 项目简介

Zebra-CICD 是一个基于 Go 语言开发的持续集成与持续部署（CI/CD）平台，旨在简化应用程序的构建、测试和部署流程。该项目采用了现代化的技术栈，包括 Gin 框架、GORM ORM、PostgreSQL 数据库以及 Kubernetes 集群管理。

## 核心功能

- 项目配置存储：支持项目的元数据管理。
- 发布任务存储：记录并跟踪部署任务的状态。
- GitLab / Harbor 客户端：提供与 GitLab 和 Harbor 的交互能力。
- 异步任务处理：模拟镜像构建与部署过程。
- 结构化日志：使用 Zap 实现全局日志记录。
- API 文档：通过 Swagger UI 提供直观的接口文档。

## 技术栈

- 后端框架：Gin
- 数据库：PostgreSQL（通过 GORM 操作）
- 日志系统：Zap + Lumberjack
- 配置管理：Viper
- 外部集成：GitLab、Harbor、Jenkins、Kubernetes
- API 文档：Swagger UI

## 目录结构
```text
zebra-cicd/
├── config/           # 配置文件管理
│   ├── config.go
│   └── configs.yaml
├── docs/             # API 文档
├── internal/         # 核心业务逻辑
│   ├── api/          # API 控制器
│   ├── core/         # 核心服务（GitLab/Harbor/Jenkins/K8s 客户端）
│   ├── handler/      # 数据库操作（CRUD）
│   ├── model/        # 数据库模型
│   ├── service/      # 业务逻辑编排层
│   └── types/        # 公共类型定义
├── pkg/              # 通用工具包
│   ├── log/          # 日志模块
│   ├── middleware/   # 中间件
│   ├── ssh/          # SSH 客户端
│   └── timeutil/     # 时间工具
├── scripts/          # 启动和构建脚本
├── main.go           # 应用入口
└── go.mod            # Go 模块依赖

```


## 快速开始

1. 准备 PostgreSQL，并设置环境变量：
    - ZEBRA_DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
    - ZEBRA_GITLAB_TOKEN=your_token
    - ZEBRA_GITLAB_URL=https://gitlab.com
    - ZEBRA_HARBOR_URL=https://harbor.example.com
    - ZEBRA_PORT=9527

2. 下载依赖并运行：
   - `go mod tidy`
   - `go run main.go`

3. 打开接口文档（Swagger UI）：
   - 访问地址：http://127.0.0.1:9527/docs
   - 安装 `go get github.com/gin-contrib/cors` 实现自定义cors
   - 安装 `go get github.com/swaggo/gin-swagger` 
   - 安装 `go get github.com/swaggo/files` 实现swagger动态生成
   - `swag init -g main.go`


### 对接k8s集群
- k8s集群权限创建
    ```yaml
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: zebra-sa
      namespace: default
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: zebra-cluster-role
    rules:
    - apiGroups: [""]
      resources: ["nodes", "pods", "services", "namespaces", "configmaps", "secrets", "events",  "jobs", "cronjobs"]
      verbs:  ["create", "get", "list", "watch", "update", "patch", "delete"]
    - apiGroups: ["apps"]
      resources: ["deployments", "statefulsets", "daemonsets"]
      verbs: ["create", "get", "list", "watch", "update", "patch", "delete"]
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: zebra-cluster-binding
    subjects:
    - kind: ServiceAccount
      name: zebra-sa
      namespace: default
    roleRef:
      kind: ClusterRole
      name: zebra-cluster-role
      apiGroup: rbac.authorization.k8s.io
    ```

- 拿到权限token
  - 获取 Service Account 的 Secret 名称
    SECRET_NAME=$(kubectl get serviceaccount zebra-sa -o jsonpath='{.secrets[0].name}')
  - kubectl get secrets
    kubectl describe secret <default-token-name>
    kubectl get secret <default-token-name> -o jsonpath='{.data.token}' | base64 -d
- 获取 Token
  - kubectl get secret $SECRET_NAME -o jsonpath='{.data.token}' | base64 -d

