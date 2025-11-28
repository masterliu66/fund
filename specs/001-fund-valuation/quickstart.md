# Quickstart: 基金估值分析应用

## 前置条件

- 已安装 Go 1.x
- 已安装 MySQL，并初始化所需数据库与表
- 已配置数据库连接（例如在配置文件或环境变量中设置 DSN）
- 本仓库已获取到本地：`git clone ...`

## 启动步骤

1. 安装依赖（如果使用 Go modules）：

   ```bash
   go mod tidy
   ```

2. 确认数据库可用并已迁移：

   - 创建 `fund` 相关库和表：
     - `fund`（基础信息表）
     - `fund_nav_history`（历史净值表）
   - 保证应用的 `dao.Db` 能连通 MySQL。

3. 启动后端服务：

   ```bash
   go run main.go
   ```

   - 服务默认监听：`http://localhost:8000`
   - Gin 路由通过 `web.NewRouter()` 注册：
     - `GET /fund`
     - `GET /funds/{fund}`
     - `POST /funds/{fund}`

4. 访问前端页面

   - Gin 使用 `template` 渲染 `frontend/index.html`（或对应模板）  
   - 在浏览器打开：

     ```text
     http://localhost:8000/
     ```

   - 页面功能：
     - `/fund` 接口加载预设基金列表及各指标
     - `/funds/{code}` 接口根据输入代码查询单只基金指标
     - 列表中可查看 TP 分位、历史高低点、当日估值、涨幅等
     - 新增列展示估值状态、投资建议、预期年化收益

5. 定时任务（可选）

   - 应用中通过 `service.StartCron` 启动定时任务：
     - 如：`insertFunds` 用于自动补录或刷新基金历史数据
   - 若在开发环境不需要定时任务，可暂时注释相关调用。

## 验证

- 打开浏览器，输入已有基金代码，验证：
  - 能正常拉取历史高/低点、TP 分位和当日估值；
  - 能看到「估值状态」「投资建议」「预期年化收益」三列；
  - 输入无效代码时，前端有友好的错误提示。
