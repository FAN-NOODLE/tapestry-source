# gw-scrm-go


#### TODO

- [X] 配置文件二进制 [go-bindata](https://github.com/go-bindata/go-bindata)
- [x] 基于swagger文档
- [x] migrate方案
- [x] CI/CD流程
- [ ] 示例接口

#### 目录结构
```
.
├── Makefile
├── README.md
├── ci-cd                           ci/cd 配置
│   ├── Jenkinsfile.groovy          
│   ├── README.md
│   ├── branch-down.groovy          分支下线脚本
│   ├── build                       打包是使用的应用配置和脚本
│   │   ├── Dockerfile
│   │   ├── config.dev.yaml
│   │   └── config.local.yaml
│   └── deployments                 各个环境的发布配置
│       ├── k8s-dev
│       │   └── demo.yaml
│       ├── k8s-prod
│       └── k8s-test
├── cmd                         
│   └── http
│       ├── api.go
│       ├── endpoints.go
│       ├── main.go
│       └── main_test.go
├── db                              数据库相关脚本
│   ├── migrations                  数据库版本管理脚本
│   └── templates                   程序中使用到的查询 sql
│       └── t_demo                  表名或模块名。
│           ├── insert.sqltpl       命名 sql。每个 sql 一个文件
│           ├── select-id.sqltpl
│           └── select.sqltpl
├── docs                            文档说明
│   ├── README.md
│   └── third                       三方 sdk 使用说明
│       └── README.md
├── git-ci                          git 命令封装，做代码静态检测
├── go-private                      go 私库拉取命令
├── go.mod
├── go.sum
├── internal
│   ├── assets
│   │   └── assets.go
│   ├── config                      应用程序配置
│   │   ├── config.go
│   │   └── struct.go
│   ├── modules
│   │   └── demo
│   │       ├── delivery            对外发布接口
│   │       │   ├── http            http 方式发布
│   │       │   │   ├── handle.go   http 发布方式具体实现封装
│   │       │   │   └── route.go    返回 http 实现是路由描述，初始化领域服务，选择具体的实现 respository
│   │       │   └── translator.go   领域实现和传输转换类
│   │       ├── init.go             提供领域服务初始化方式，供 delivery 调用，注入 具体实现的 respository
│   │       ├── entity.go           实体
│   │       ├── mocks               测试 mock
│   │       ├── repository.go       存储接口
│   │       ├── respository         存储具体实现
│   │       │   └── mysql
│   │       │       ├── impl.go
│   │       │       └── init.go
│   │       └── usercase.go         应用服务
│   └── pkg                         内部类库
│       ├── cache
│       │   ├── cache.go
│       │   └── goodsCache.go
│       ├── db
│       │   └── init.go
│       ├── middleware
│       │   └── access.go
│       ├── mysqlRepositoyBase.go   
│       ├── redis
│       │   └── init.go
│       └── tracing
│           └── skywalking.go
├── pkg                         本服务对外提供的 sdk
│   └── README.md
├── scripts                     make 命令使用的脚本，开发过程中使用的脚本
└── sonar-project.properties    sonar 工程配置

```

#### quick start

```shell script
    # 更改库获取方式
    git config --global url."ssh://git@gitlab-ce.k8s.tools.vchangyi.com:32201".insteadOf "https://gitlab-ce.k8s.tools.vchangyi.com"
    git config --global url."ssh://git@gitlab-ce.k8s.tools.vchangyi.com:32201".insteadOf "http://gitlab-ce.k8s.tools.vchangyi.com"
    # 设置私库。如果要永久有效请配置到环境变量
    export GOPRIVATE=gitlab-ce.k8s.tools.vchangyi.com
    # 下载依赖
    go mod download
    # 运行
    make run M=http E=local
    # 编译
    go build .
```

### 数据库升级

依赖 docker 运行

```shell script
# 创建一个合并版本文件
make migrate-create V=1.0.0 DESC=init

# 升级或回退到特定版本
make migrate V=1.0.0 DB='mysql://root:toor@tcp\(127.0.0.1:3306\)/test?multiStatements=true'
make migrate V=1.1.0 DB='mysql://root:toor@tcp\(127.0.0.1:3306\)/test?multiStatements=true'
make migrate V=1.1.3 DB='mysql://root:toor@tcp\(127.0.0.1:3306\)/test?multiStatements=true'

# 强制设置当前版本号，注意：不会执行合并或回退操作
make migrate-fix V=1.1.3 DB='mysql://root:toor@tcp\(127.0.0.1:3306\)/test?multiStatements=true'

# 显示帮助信息
make migrate-help 

```

