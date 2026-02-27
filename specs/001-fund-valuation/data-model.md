# Data Model: 基金估值分析

## Entity: Fund

- **Purpose**: 表示单只基金的基础信息
- **Fields**:
  - `code` (string)  
    - 基金代码，主键或唯一索引
  - `name` (string)  
    - 基金名称
  - `type` (string, optional)  
    - 基金类型（如 equity/bond/mixed/FOF 等）
- **Relationships**:
  - 1 : N → `FundNavHistory`
- **Notes**:
  - 代码应与数据源使用的编码保持一致

---

## Entity: FundNavHistory

- **Purpose**: 存储单只基金每日净值 & 估算值，用于计算历史高低点和分位
- **Fields**:
  - `id` (auto-increment)  
  - `fund_code` (string, FK → Fund.code)  
  - `date` (date)  
  - `nav` (decimal)  
    - 当日单位净值
  - `estimated_nav` (decimal, optional)  
    - 当日估算净值（如有）
  - `return_rate` (decimal, optional)  
    - 当日涨跌幅（%）
- **Indexes**:
  - `(fund_code, date)` 组合唯一索引
- **Notes**:
  - 通过该表计算：
    - 历史最高/最低
    - 近 1 年、90 天、上月、本月各区间的高低点
    - TP20/TP80 等分位

---

## Entity: FundInfoReport (View/DTO)

- **Purpose**: 服务接口 `/fund` 与 `/funds/{fund}` 的聚合视图/DTO，聚合某只基金当前相关指标和分析结论
- **Fields**（对应 Go struct）:
  - `fundCode` (string)
  - `name` (string)
  - 区间高低：
    - `historyMaxDwjz`, `historyMinDwjz`, `historyAvgDwjz`
    - `lastYearMaxDwjz`, `lastYearMinDwjz`
    - `lastSeasonMaxDwjz`, `lastSeasonMinDwjz`
    - `lastMonthMaxDwjz`, `lastMonthMinDwjz`
    - `maxDwjz`, `minDwjz`（当前周期）
  - 分位：
    - `tp80MinDwjz`, `tp80MaxDwjz`
    - `tp85MinDwjz`, `tp85MaxDwjz`
  - 当日：
    - `gsz`（当日估值）
    - `gszzlFormat`（当日涨幅字符串）
  - 估值与投资结论：
    - `valuationStatus` (`"高估" | "低估" | "合理" | "未知"`)
    - `valuationScore` (number)
    - `investAdvice` (string)
    - `riskNote` (string)
  - 预期收益：
    - `expectedReturnMin` (number, 0.08 = 8%)
    - `expectedReturnMax` (number)
    - `expectedReturnNote` (string)
  - 数据状态：
    - `dataStatus` (`"OK" | "PARTIAL" | "MISSING"`)
    - `dataStatusNote` (string)
- **Notes**:
  - 可以是数据库视图，也可以纯 Go 聚合 DTO，不强制落表

---

## Entity: FundConfig

- **Purpose**: 管理系统中展示的基金配置信息，决定哪些基金会在前端列表中显示
- **Fields**:
  - `id` (int64, auto-increment)  
    - 主键ID
  - `fund_code` (string, unique)  
    - 基金代码，6位数字
  - `fund_name` (string)  
    - 基金名称
  - `fund_type` (int32)  
    - 基金类型（0:普通基金, 1:海外基金）
  - `created_at` (timestamp)  
    - 创建时间
  - `updated_at` (timestamp)  
    - 更新时间
  - `record_version` (int32)  
    - 记录版本号
  - `is_deleted` (int32)  
    - 软删除标记（0:正常, 1:已删除）
- **Indexes**:
  - `(fund_code)` 唯一索引
  - `(fund_type, is_deleted)` 复合索引
- **Relationships**:
  - 1 : N → `fund_info` (通过fund_code关联)
- **Notes**:
  - 软删除机制保证数据完整性
  - 基金类型使用枚举值，便于前端展示和筛选
  - 记录版本号用于并发控制

---

## Entity: UserWatchlist (如未来需要)

- **Purpose**: 保存用户自选基金列表
- **Fields**:
  - `id` (auto-increment)
  - `user_id` (string or int)
  - `fund_code` (string, FK → Fund.code)
  - `created_at` (timestamp)
- **Notes**:
  - 当前实现里可能还未涉及，可作为后续扩展
