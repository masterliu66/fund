# Implementation Plan: 基金估值分析应用

**Branch**: `001-fund-valuation` | **Date**: 2025-11-28 | **Spec**: `specs/001-fund-valuation/spec.md`
**Input**: Feature specification from `/specs/001-fund-valuation/spec.md`

## Summary

为指定基金代码收集历史净值与当日估值数据，计算历史高低点、TP 分位等指标，在前端展示单基金与多基金的估值概览、当前高估/低估状态、是否适合投资以及预期收益区间。

## Technical Context

**Language/Version**: Go 1.x (Golang)  
**Primary Dependencies**: Gin (HTTP web framework), sqlx (DB access), standard library (net/http, time, etc.)  
**Storage**: MySQL (primary OLTP database for 基金净值与历史数据)  
**Testing**: [NEEDS CLARIFICATION: 单元测试框架选择（例如 testing + testify？）以及是否需要集成测试方案]  
**Target Platform**: Linux server (backend), modern desktop browsers (frontend)  
**Project Type**: web  
**Performance Goals**:  
- 后端接口 P95 延迟 < 300ms（在正常负载下，单次基金查询 & /fund 列表查询）  
- 支持同时在线查询基金数量：至少数百只基金列表的刷新  
**Constraints**:  
- P95 接口响应时间 < 300ms  
- 单次接口查询中数据库访问次数可控（避免 N+1 查询）  
- 页面首次加载体积适中（纯原生 JS，无大型前端框架）  
**Scale/Scope**:  
- 基金代码规模：数十到数百只自选基金  
- 单个基金历史净值记录：可达数千条  
- 单实例支撑的日活用户量：数十到数百（个人/小团队使用场景）

## Constitution Check

- Code Quality & Readability  
  - Go 代码保持清晰包结构：`dao/`, `model/`, `service/`, `web/`, `frontend/` 等分层职责明确。  
  - 公共 struct（如 FundInfoReport）和 handler（如 GetFundInfo）在含义不直观时补充简要注释。

- Testing Discipline  
  - 计划为关键策略计算逻辑（基金高/低估判断、投资建议生成）增加单元测试。  
  - `/fund` 与 `/funds/{fund}` 路由增加最基本的 handler 层测试或集成测试（至少覆盖成功/错误路径）。  
  - 具体测试框架仍待确定（见 Technical Context 中的 NEEDS CLARIFICATION）。

- User Experience Consistency  
  - 前端列表页使用统一的表格布局与列显示/隐藏逻辑（已有 columnStates 机制）。  
  - 查询失败、代码无效、数据缺失等情况提供统一样式的错误提示文案。  
  - 为新增加的“估值状态 / 投资建议 / 预期收益”列设计统一的展示格式。

- Performance & Resource Efficiency  
  - `/fund` 接口对基金列表查询应避免 N+1 查询，优先使用批量查询与合适索引。  
  - 对大列表场景可考虑限制一次返回数量或增加简单分页/过滤（如果未来规模增大）。  
  - 定时任务 `insertFunds` 避免在高峰时间对数据库造成过大压力。

- Operational Reliability & Observability  
  - 接口失败路径记录至少包含 fund code、请求路径和错误信息（不泄露敏感信息）。  
  - 后端错误日志使用统一格式，方便后续查询和告警接入。  
  - 可以考虑为关键接口添加简单指标（调用次数、错误数），后续接入监控系统。

## Project Structure

```text
specs/001-fund-valuation/
├── spec.md
├── plan.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── checklists/
    └── requirements.md
```

**Structure Decision**: 使用单仓库 + 单 Go 服务 + 原生 JS 前端的结构，后端继续沿用当前 `dao/model/service/web/frontend` 分层。

## Complexity Tracking

> 如未来引入更复杂的估值模型或缓存/队列等，在此处记录复杂度来源与简化方案评估。
