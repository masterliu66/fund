---

description: "Task list template for feature implementation: 基金估值分析应用"
---

# Tasks: 基金估值分析应用

**Input**: Design documents from `/specs/001-fund-valuation/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), data-model.md, contracts/, quickstart.md

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 确认项目结构与依赖、文档和环境准备

- [ ] T001 更新 plan.md 中 Go 版本与实际环境保持一致（已完成可在此勾选）
- [ ] T002 确认 MySQL 数据库结构与 `fund_info.sql` / `fund_record.sql` 一致，并在目标环境中执行建表脚本
- [ ] T003 [P] 在 `README_API_CHANGES.md` 或 docs 中补充本次特性简介和访问方式
- [ ] T004 [P] 在 `specs/001-fund-valuation/quickstart.md` 中补充实际数据库连接配置示例（如 DSN 模板）

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 为所有用户故事准备共享的数据访问与基础结构

- [ ] T005 实际检查并完善 DAO 层：在 `dao/fund_dao.go` / `dao/fund_record_dao.go` 中确保存在读取历史净值和当日估值所需的查询
- [ ] T006 在 `model/fund_info_report.go` 中为 `FundInfoReport` 补充估值状态、投资建议、预期收益、数据状态等新字段（仅定义 struct 与 json tag）
- [ ] T007 在 `service` 层（如 `service/fund_service.go` 或等效文件）中集中封装查询 `FundInfoReport` 所需的数据聚合逻辑（避免控制器直接拼装）
- [ ] T008 确保 `web.NewRouter()` 中正确注册 `/fund` 和 `/funds/{fund}` 路由，并返回 `FundInfoReport` 的 JSON 字段与前端约定一致

**Checkpoint**: DAO 和基础 service 能返回完整的 FundInfoReport 基础字段（尚未包含估值结论字段也可以），接口可正确访问。

---

## Phase 3: User Story 1 - 查询单只基金估值概览 (Priority: P1) 🎯 MVP

**Goal**: 用户输入基金代码后，可以看到单只基金的完整估值概览（包括历史高低点、分位、当日估值和基本结论）。

**Independent Test**: 仅实现 /funds/{fund} 相关逻辑时，用户仍可通过输入代码，完成从查询到阅读结论的完整闭环。

### Tests for User Story 1 (可选，但推荐)

- [ ] T009 [P] [US1] 在 `service` 层添加单元测试文件（如 `service/fund_strategy_test.go`），覆盖高估/低估/合理三种典型场景
- [ ] T010 [P] [US1] 在 `web` 层添加最基本的 handler 测试或集成测试（如 `web/fund_controller_test.go`），验证 `/funds/{fund}` 成功与错误路径

### Implementation for User Story 1

- [ ] T011 [P] [US1] 在 `model/fund_info_report.go` 中根据最终 JSON 设计，确认并补充所有需要的字段（如 `ValuationStatus`, `InvestAdvice`, `ExpectedReturnMin/Max` 等）
- [ ] T012 [P] [US1] 在 `service` 包中新增 `FillValuationFields` 辅助函数（例如 `service/valuation_helper.go`），根据 Gsz 与 TP80 区间填充估值状态与建议
- [ ] T013 [US1] 在 `service.CalFundsStrategy` / `service.CalFundsStrategy2` 中调用 `FillValuationFields`，确保返回的 `FundInfoReport` 包含估值状态与预期收益字段
- [ ] T014 [US1] 确认 `web.GetFundInfo` 返回的 JSON 中包含所有新字段（检查 gin 的 JSON 序列化结果与前端需要的字段名一致）
- [ ] T015 [P] [US1] 在 `frontend/index.html` 的 `search()` 与 `show()` 逻辑中，读取 `/funds/{code}` 返回的数据并在页面某处展示单只基金的估值状态与投资建议（如弹出或详情区域）
- [ ] T016 [US1] 为单只基金查询失败场景添加统一错误提示（`msg("请输入代码")` 之外，如无效代码、接口错误等），确保 UX 一致

**Checkpoint**: 用户可以输入单个基金代码，看到历史高低点、TP 指标、当日估值和简洁的估值结论，不依赖多基金对比功能。

---

## Phase 4: User Story 2 - 对比多只基金当前估值与预期收益 (Priority: P2)

**Goal**: 用户可以在列表中同时查看多只基金的估值状态与预期收益区间，用于筛选和排序。

**Independent Test**: 即使只实现多基金对比视图（在 US1 已完成前提下），用户也能通过添加多只基金、查看对比结论完成独立决策过程。

### Tests for User Story 2 (可选)

- [ ] T017 [P] [US2] 为 `/fund` 接口添加集成测试（例如在 `web` 或 `httpt` 目录中），验证返回的列表长度、字段完整性以及错误处理

### Implementation for User Story 2

- [ ] T018 [P] [US2] 在 `service.CalFundsStrategy` / `CalFundsStrategy2` 中确认批量基金时也应用估值和预期收益逻辑（循环内调用 `FillValuationFields`）
- [ ] T019 [P] [US2] 在 `frontend/index.html` 表头中增加「估值状态」「投资建议」「预期年化收益」三列，并在 `DEFAULT_COLUMN_STATES` 中添加对应布尔值
- [ ] T020 [US2] 在 `frontend/index.html` 的 `show()` 函数中，按列顺序为每一行渲染新字段（`valuationStatus`, `investAdvice`, `expectedReturnMin/Max`），并格式化为百分比区间
- [ ] T021 [US2] 确保前端列显示/隐藏逻辑（`columnStates`、复选框控制）支持新列，不破坏现有列行为
- [ ] T022 [US2] 在前端增加简单排序/过滤交互（如根据估值状态或预期收益大致筛选），若暂不实现可在 tasks 中记录为后续优化

**Checkpoint**: `/fund` 列表中每只基金均显示估值状态、投资建议和预期收益区间，用户可以快速对多只基金进行初步筛选。

---

## Phase 5: User Story 3 - 评估当前是否适合买入/加仓 (Priority: P3)

**Goal**: 用户可以看到针对某只基金的「是否适合当前投资」结论和主要依据说明。

**Independent Test**: 在已有估值计算基础上，只实现投资建议文案与解释模块时，用户也可以单独访问某只基金详情并理解建议原因。

### Tests for User Story 3 (可选)

- [ ] T023 [P] [US3] 为 `FillValuationFields` 或对应策略函数增加更多边界测试用例（如极端高估、极端低估、数据缺失）

### Implementation for User Story 3

- [ ] T024 [US3] 在 `service` 层扩展投资建议规则（如根据估值状态 + 历史波动情况调整 `InvestAdvice` 与 `RiskNote`），保证文案对普通用户友好
- [ ] T025 [US3] 在前端为单只基金详情或弹出区域展示「是否适合投资」结论与 1–3 条要点依据（可使用现有 index.html 或新增简单详情视图）
- [ ] T026 [US3] 针对数据不足或模型不适用场景，在后端统一设置 `ValuationStatus = "未知"` 并返回明确原因说明，前端配套展示「暂无法给出可靠建议」提示

**Checkpoint**: 用户可以在查看单只基金时，获得直观的「适合/谨慎/观望」类建议以及简短解释。

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 跨故事的优化和质量提升

- [ ] T027 [P] 文档更新：在 `docs/` 或 `README_API_CHANGES.md` 中补充接口字段变更（特别是 `/fund` 与 `/funds/{fund}` 的新字段）
- [ ] T028 [P] 代码清理：整理 `service` / `web` / `model` 中与基金估值相关的函数命名与注释，使之与 spec/plan 中术语一致
- [ ] T029 [P] 性能检查：对 `/fund` 与 `/funds/{fund}` 做一次简单压测，确认在预计数据规模下 P95 延迟 < 300ms，并在 `specs/001-fund-valuation/research.md` 或 notes 中记录结果（如果创建该文件）
- [ ] T030 [P] 日志与错误处理：检查 `web` 和 `service` 层对于主要错误路径是否有足够日志（含 fund code 和错误原因），并在必要处补充
- [ ] T031 校对规格与实现：对照 `specs/001-fund-valuation/spec.md` 中的 FR 和 SC 条目，逐一验证是否已有对应实现与测试/验证方案

---

## Dependencies & Execution Order

- **Phase 1**: 无依赖，可立即开始
- **Phase 2**: 依赖 Phase 1 完成（环境和结构确认）
- **Phase 3 (US1)**: 依赖 Phase 2 完成
- **Phase 4 (US2)**: 建议在 US1 完成后进行，但多数任务可与 US1 部分并行（如前端列结构调整）
- **Phase 5 (US3)**: 依赖 US1 的估值计算基础完成
- **Phase 6**: 在主要用户故事完成后进行

---

## Parallel Opportunities

- DAO 与 model 字段补充（T005, T006）可以并行
- 前端列头/列渲染调整（T019, T020, T021）可以与部分后端逻辑并行
- 测试用例编写（T009, T010, T017, T023）可以与实现交错进行

---

## Implementation Strategy

- **MVP（US1）**: 优先完成单只基金估值概览功能（Phase 3），确保用户可以输入代码并看到完整的估值信息与结论
- **增量迭代**: 依次实现多基金对比（US2），再实现更丰富的投资建议与解释（US3）
- **收尾与优化**: 最后统一处理文档、性能、日志与清理（Phase 6），确保满足宪章中的质量和性能要求
