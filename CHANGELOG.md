# Changelog

所有项目的显著变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [v0.5.0] - 2025-03-25

- :bookmark: Bump version to 0.5.0 (98f9b00)
- :recycle: 重构日志收集和上传逻辑 (fa44cde)
- :page_facing_up: Create LICENSE (9da9e9e)
- :memo: 新增架构设计图 (46c0d6c)
- :memo: 更新自述文件 (5a08e5f)


## [v0.4.1] - 2025-03-20

- :bookmark: Bump version to 0.4.1 (71a64fb)
- :recycle: 重构 Makefile 以使用环境变量存储配置 URL，并更新发布工作流以利用 secrets 保存上传和下载 URL (e609579)
- :pencil2: 更新安装脚本和文档中的下载链接，替换为示例域名 (9bd14ef)
- :white_check_mark: 添加单元测试，覆盖命令行执行、时间参数解析、收集器功能及处理结果验证，增强代码的测试覆盖率和稳定性 (e9801c2)


## [v0.4.0] - 2025-03-13

- :memo: 更新 CHANGELOG.md，添加 v0.4.0 版本变更日志 (4b0202f)
- :bookmark: Bump version to 0.4.0 (c93f08c)
- :bug: 修复无法正确处理多行日志条目的问题 (abbd5d2)
- :sparkles: 添加新的日志处理器支持，包括 VisionLogViewer 和 RobotDriverNode，更新工厂模式以支持新处理器类型 (4e303ab)
- :recycle: 更新支持程序列表获取方式，使用工厂模式替代直接引用 (d2d71c8)


## [v0.3.0] - 2025-03-12

- :memo: 更新 CHANGELOG.md，添加 v0.3.0 版本变更日志 (3414e81)
- :bookmark: Bump version to 0.3.0 (a22492c)
- :sparkles: 更新补全脚本，添加 supported-programs 命令支持 (b0ac538)
- :pencil2: 更新补全脚本中的程序名称 xyz-hmi-server 为 xyz-max-hmi-server (9c36fee)
- :sparkles: 新增命令行自动补全功能，支持Bash/Zsh/Fish/PowerShell (2adae5b)
- :sparkles: 新增 supported-programs 命令，用于展示当前支持日志收集程序 (566eb22)
- :recycle: 更新 ProcessorType (ad791a9)


## [v0.2.0] - 2025-03-12

- :memo: 更新 CHANGELOG.md (733d52b)
- :bookmark: Bump version to 0.2.0 (822de7b)
- :white_check_mark: 新增日志处理器二进制打包模块的单元测试 (39f674d)
- :white_check_mark: 新增 JSON 文件处理器的单元测试 (6f96630)
- :test: 添加多个模块的单元测试 (3bdf5c3)
- :test: 完善日志处理内容测试用例 (13c7f76)
- :recycle: 重构日志处理器文件处理逻辑 (86756be)
- :memo: 更新 FilterFiles 函数文档 (bd69bf0)
- :test: 添加用户操作日志处理器的单元测试 (472228a)
- :test: 更新测试用例以适应新的日志处理接口 (6b9a023)
- :recycle: 重构日志处理器模块 (737a928)
- :recycle:  重构 processor (51e09bd)
- :lock: 为安装脚本添加随机参数以防止缓存 (0f5e8f3)
- :mute: 移除调试信息 (c18f9e2)
- :sparkles: 添加变更日志生成脚本以及初始化变更日志 (6bb85b0)


## [v0.1.1] - 2025-03-06

- :bookmark: Bump version to 0.1.1 (07955c9)
- :sparkles: 支持通过构建时注入配置 URL (9551b64)
- :bug: 改进安装脚本的版本和下载链接获取逻辑 (8bca80d)
- :rocket: 更新安装脚本以支持动态版本和下载链接 (6498f6f)
- :recycle: 重构远程配置管理逻辑 (b16dc52)
- :sparkles: 添加配置文件示例和更新 .gitignore (a837f55)
- :recycle: 重构远程配置 URL 获取方法 (524c495)
- :fire: 移除 GitHub Release 创建步骤 (7566241)


## [v0.1.0] - 2025-03-05

- :sparkles: 完善手动触发版本发布工作流 (4148698)
- :sparkles: 支持手动触发版本发布 (1836bd2)
- :lock: 更新远程配置链接 (75bc55b)
- :rocket: 添加 GitHub Actions 自动发布工作流 (899ed0b)
- :rocket: 添加自动发布脚本 (f4d6df5)
- :zap: 优化 ZipDirectory 函数的并发压缩性能 (34922d5)
- :zap: 优化日志收集器的并发处理性能 (79dc4df)
- :recycle: 重构文件过滤和排序逻辑 (07b5542)
- :sparkles: 增强日志处理器的文件匹配和路径处理 (fa03d30)
- :bug: 改进 StudioMax 日志处理器的错误处理 (a3bd8f1)
- :sparkles: 新增程序日志选择功能 (c719074)
- :arrow_up: 清理依赖并移除未使用的模块 (c697ff6)
- :art: 重构日志格式化器，增强日志输出样式 (07f259e)
- :sparkles: 增强日志收集器功能和工厂模式支持 (bc4e04e)
- :recycle: 调整 service.go 中的目录路径和导入路径 (1db67ca)
- :see_no_evil: 更新 .gitignore，忽略更多测试结果目录 (0dc6930)
- :white_check_mark: 新增 StudioMax 日志处理器的测试用例 (07ff178)
- :sparkles: 新增 StudioMax 日志处理器 (b7c6f72)
- :recycle: 重命名 hmiserver 为 hmi_server (3b1ab03)


## [v0.0.3] - 2025-03-04

- :bookmark: 引入版本管理和构建时信息注入 (f32d7bc)
- :white_check_mark: 新增 HMI 日志处理器测试用例 (89c5d2b)
- :sparkles: 优化 HMI 日志处理和文件过滤逻辑 (2e73b85)
- :recycle: 重命名 (39230a2)


## [v0.0.2] - 2025-03-04

- :bookmark: 更新版本为 0.0.2 (8ecbe4c)
- :construction_worker: 新增 Makefile 打包目标平台二进制文件的功能 (009338b)
- 📝 更新 README.md 安装和使用说明 (b6b42ff)
- :recycle: 重命名日志文件保留选项，提高配置语义性 (9ec48a8)
- 🔧 调整日志收集完成后的日志输出逻辑 (41d512e)
- ✨ 优化日志收集后处理逻辑 (27bd353)
- 🔧 简化命令行上传配置 (6d8376a)
- 🔧 调整 Cloudreve 分享文件链接参数 (d1d9abb)
- ✨ 为更新命令添加 TUI 界面和自动更新支持 (ecd8689)
- ✨ 为版本命令添加简单模式和 TUI 界面支持 (eb2f66a)
- :recycle: 重构命令行模块，拆分功能到独立文件 (af0370e)
- :coffin: 移除 UI 相关模块，简化命令行交互逻辑 (b05f6f3)
- 🔧 优化 HMI 服务器日志文件名处理逻辑 (80d26fa)
- :construction_worker: 优化 Makefile 构建目录结构 (66d23af)
- :memo: 简化 README 中的日志收集命令示例 (3828f48)
- 🔧 增强 Makefile，支持跨平台二进制构建 (89e0e86)
- 🔧 Update go.mod dependencies (0fb3a0d)
- ✨ 新增日志过滤组件的单元测试和增强处理逻辑 (7cb80cd)
- 🔧 重构日志收集和处理组件，引入 logrus 日志库 (1539a1f)
- ✨ 新增程序更新和版本管理功能 (54fcf8e)
- 🔧 优化 Makefile，增加二进制文件构建选项 (a0ec80e)
- 🔧 支持开发和生产环境的远程配置URL (51c811c)
- ✨ 新增 Windows 和 Linux 安装脚本 (97b6c55)
- ✨ 重构上传组件，支持多云存储提供商 (e941727)
- 🔧 重构表单和日志组件，优化交互细节 (7bd3bef)
- ✨ 重构版本更新逻辑，支持简单模式和 TUI 模式的版本检查 (9c09794)
- ✨ 新增版本更新和远程配置管理功能 (c1305a6)
- 🔧 更新依赖并清理无用的第三方库 (2d1d75e)
- 🔧 移除日志配置保存输入组件 (5651110)
- 🔧 精简日志过滤和表单交互逻辑 (907ee37)
- ✨ 优化表单交互和日志管理机制 (bc47796)
- ✨ 重构项目架构，引入交互式UI和服务层 (998f05a)
- :sparkles: 支持 xyz-hmi 日志的收集 (2b7bc7e)
- 优化日志收集和处理逻辑 (340f792)
- 初始化提交 (cd4d1ca)


