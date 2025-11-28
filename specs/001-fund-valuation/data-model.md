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

## Entity: UserWatchlist (如未来需要)

- **Purpose**: 保存用户自选基金列表
- **Fields**:
  - `id` (auto-increment)
  - `user_id` (string or int)
  - `fund_code` (string, FK → Fund.code)
  - `created_at` (timestamp)
- **Notes**:
  - 当前实现里可能还未涉及，可作为后续扩展
