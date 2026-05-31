# 学生成绩管理系统

基于 Wails 框架开发的桌面端学生成绩管理系统，Go 后端 + Vue 3 前端。

## 功能概览

### 学生管理
- 学生信息的增删查改
- 支持按姓名搜索、分页浏览
- 管理员专属：新增和删除学生

### 课程管理
- 课程信息的增删查改
- 按学期、课程名称管理
- 老师只能删除自己创建的课程

### 成绩管理
- 成绩录入、修改、删除（含绩点自动计算）
- 多条件查询：按学号、姓名、课程、学期
- 批量导入：粘贴 CSV 文本一键导入，逐行错误提示
- 批量调整：按课程 + 分数段批量加减分
- 数据汇总：跨课程、跨学期成绩汇总
- 成绩单导出：生成标准 Excel 成绩单

### 绩点计算
- 公式：`score ≥ 60 → gp = round(score/10 − 5, 1)`，`score < 60 → gp = 0`
- 绩点规则可查看和编辑（存储在 `data/gpa_rules.txt`）
- 支持批量重算全部学生绩点

### 统计分析
- 学生统计：平均分、GPA、总学分、课程数
- 综合排名：按绩点降序排名
- 课程报表：选择课程导出 Excel（含平均分、及格率、分数段分布）
- 学生报表：选择学生导出 Excel（含各科明细、总分、绩点、排名）

### 数据管理
- 按学期或课程备份成绩数据
- 从备份目录恢复数据
- 操作日志：所有成绩操作全记录，支持多条件追溯
- 日志导出 Excel

### 权限管理
- 管理员（admin）：全部功能，含用户管理、操作日志
- 老师（数字工号）：学生搜索、课程管理、成绩录入/修改/删除（仅自己录入的）

### 数据校验
- 分数 0-100 范围校验
- 学号/课程存在性校验
- 校验失败生成错误日志 `data/error.log`

## 技术栈

| 层 | 技术 |
|---|---|
| 框架 | Wails v2 |
| 后端 | Go + GORM + SQLite |
| 前端 | Vue 3 + Vite + Pinia + Vue Router |
| Excel | excelize |
| 密码 | bcrypt |

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- Wails CLI（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）

### 运行开发模式
```bash
wails dev
```

### 构建可执行文件
```bash
wails build
cp build/bin/Student-Grade-Management-System.exe .
```
产物 `Student-Grade-Management-System.exe` 在项目根目录，与 `database/`、`data/` 同级，双击即可运行。

### 构建安装包
需安装 [NSIS](https://nsis.sourceforge.io/)：
```bash
wails build -nsis
```
产物在 `build/bin/Student-Grade-Management-System-amd64-installer.exe`

### 默认账号
| 角色 | 用户名 | 密码 |
|---|---|---|
| 管理员 | admin | 12345678 |
| 老师 | 由管理员创建 | 8-12 位 |

## 项目结构

```
├── main.go              # 应用入口
├── app.go               # 应用结构
├── wails.json           # Wails 配置
├── backend/
│   ├── api/             # API 层（对前端暴露的方法）
│   ├── config/          # 数据库初始化与配置
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   │   └── csv/         # CSV 文件读写
│   ├── service/         # 业务逻辑层
│   └── utils/           # 绩点计算工具
├── frontend/
│   ├── src/
│   │   ├── views/       # 页面组件（11 个视图）
│   │   ├── components/  # 通用组件（通知弹窗、可搜索选择器）
│   │   ├── composables/ # 组合式函数（notify、clickOutside）
│   │   ├── layout/      # 主布局（侧边栏 + 内容区）
│   │   ├── router/      # 路由配置（含权限守卫）
│   │   └── store/       # Pinia 状态管理（auth）
│   └── wailsjs/         # Wails 自动生成的前端绑定
├── data/                # CSV 数据文件与规则配置
├── build/               # 构建资源（图标、NSIS 安装脚本）
├── database/            # SQLite 数据库（运行时生成，不入库）
├── export/              # Excel 导出目录（运行时生成）
└── backup/              # 数据备份目录（运行时生成）
```

## 数据存储

- **主存储**：SQLite 数据库 `database/student.db`
- **CSV 同步**：学生、课程、成绩数据同步写入 `data/` 目录下的 CSV 文件
- **成绩文件**：按 `data/grades/{学期}/{课程代码}.csv` 组织
- **操作日志**：存储在 SQLite `operation_logs` 表，可在日志界面查看和导出
- **错误日志**：校验失败记录写入 `data/error.log`
