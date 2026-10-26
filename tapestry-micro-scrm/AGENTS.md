# AGENTS.md

本文件是 `micro-scrm` 仓库的 Agent 协作入口文件，适配 Spec Kit-cn 工作流。

## 编码问题

注意本项目所有代码及文档都是utf8编码。

## 优先阅读

在处理非简单任务前，先阅读：

- `docs/CODEX_PROJECT_CONTEXT.md`

该文档包含项目上下文、模块边界、本地环境备注和长期协作规则。除非用户在当前对话中给出更新指令，否则把它视为本仓库的主要上下文来源。

## 文件编码
所有代码文件及文档均为utf8编码

## 项目概览

`micro-scrm` 是一个 Go 1.18 后端服务，主要面向企业微信 SCRM 业务。项目使用 Gin、GORM、MySQL、ClickHouse、Redis、Kafka、CMQ、Elasticsearch、腾讯云 COS 和 Azure Blob。

主要入口：

- `cmd/http/main.go`：HTTP API 服务。
- `cmd/job/main.go`：定时任务。
- `cmd/queue/main.go`：队列消费。
- `cmd/binlog-es/main.go`：binlog 同步到 Elasticsearch。
- `cmd/binlog-kafka/main.go`：binlog 同步到 Kafka。
- `cmd/binlog-ck-queue/main.go`：ClickHouse 相关 binlog 队列消费。

## 工作规则

- 只实现用户明确要求的内容。
- 不做顺手重构、格式化、重排 import 或无关清理。
- 除非需求明确要求改变行为，否则保持现有行为不变。
- 优先给出长期可维护的解决方案，不写临时兼容分支、TODO 占位或新旧实现并存代码。
- 每一处代码改动都必须能追溯到用户需求。
- 如果业务规则、API 契约、租户逻辑、权限规则、数据库结构或第三方集成语义不清楚，先向用户澄清。
- 保留工作区中已有的用户或同事改动，不要擅自回退。

## 代码组织

仓库有两块主要实现区域：

- `internal/app/`：较新的代码区域，通常采用 Router / Handler / Service / DAO / Model 分层。新增功能在任务允许时优先放这里。
- `internal/modules/`：较早的业务模块，很多代码使用 SQL 模板和旧 repository 模式。只有任务明确涉及旧模块行为时才优先修改这里。

模型、DAO、分库分表、素材库模块、Open API 等详细约定见 `docs/CODEX_PROJECT_CONTEXT.md`。

## 验证要求

处理非简单任务前，先明确：

- 要改变的具体行为。
- 可能涉及的模块和文件。
- 当前做出的假设。
- 准备运行的验证命令。

验证命令应尽量聚焦到本次改动。Go 代码改动优先运行相关包测试，例如：

```powershell
go test ./internal/app/...
```

如果验证命令无法在本地运行，需要在最终回复中明确说明。
