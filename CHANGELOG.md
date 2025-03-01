# Changelog

所有项目的显著变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [v0.4.0] - 2025-03-13

### 新特性
- :sparkles: 添加新的日志处理器支持：
  - 新增 `VisionLogViewer` 日志处理器
  - 新增 `RobotDriverNode` 日志处理器
- :sparkles: 新增通用 C++ 日志处理器，可以处理标准格式的 C++ 日志文件

### 优化
- :recycle: 重构支持程序列表获取方式：
  - 使用工厂模式替代直接引用，提高了代码的可扩展性
  - 将 `GetSupportedProcessorTypes` 方法从 `collector` 包移至 `factory` 包
  - 优化了 `ProcessorFactoryRegistry` 的结构，支持更灵活的处理器注册

### 修复
- :bug: 修复无法正确处理多行日志条目的问题：
  - 重写了 `ProcessLogContent` 方法，现在可以正确识别和处理跨多行的日志条目
  - 添加了日志条目的状态跟踪，确保多行日志条目作为一个整体被处理
  - 优化了时间范围过滤逻辑，确保多行日志条目的所有行都被正确保留或过滤

### 其他
- :bookmark: 版本号从 0.3.0 更新到 0.4.0

## [v0.3.0] - 2025-03-12

### 新特性
- :sparkles: 新增命令行自动补全功能，支持Bash/Zsh/Fish/PowerShell
- :sparkles: 新增 supported-programs 命令，用于展示当前支持日志收集程序

### 改进
- :sparkles: 更新补全脚本，添加 supported-programs 命令支持
- :pencil2: 更新补全脚本中的程序名称 xyz-hmi-server 为 xyz-max-hmi-server
- :recycle: 更新 ProcessorType

### 其他
- :bookmark: 版本升级至 0.3.0


## [v0.2.0] - 2025-03-12

### 改进
- :recycle: 重构日志处理器文件处理逻辑和模块
- :memo: 更新 FilterFiles 函数文档
- :lock: 为安装脚本添加随机参数以防止缓存
- :mute: 移除调试信息

### 测试
- :white_check_mark: 新增日志处理器二进制打包模块的单元测试
- :white_check_mark: 新增 JSON 文件处理器的单元测试
- :white_check_mark: 添加多个模块的单元测试
- :white_check_mark: 完善日志处理内容测试用例
- :white_check_mark: 添加用户操作日志处理器的单元测试
- :white_check_mark: 更新测试用例以适应新的日志处理接口

### 其他
- :bookmark: 版本升级至 0.2.0
- :sparkles: 添加变更日志生成脚本以及初始化变更日志


## [v0.1.1] - 2025-03-06

### 新特性
- :sparkles: 支持通过构建时注入配置 URL

### 改进
- :recycle: 重构远程配置管理逻辑
- :recycle: 重构远程配置 URL 获取方法
- :sparkles: 添加配置文件示例和更新 .gitignore

### 修复
- :bug: 改进安装脚本的版本和下载链接获取逻辑
- :rocket: 更新安装脚本以支持动态版本和下载链接
- :fire: 移除 GitHub Release 创建步骤

### 其他
- :bookmark: 版本升级至 0.1.1


## [v0.1.0] - 2025-03-05

### 新特性
- :sparkles: 完善手动触发版本发布工作流
- :sparkles: 支持手动触发版本发布
- :sparkles: 新增程序日志选择功能
- :sparkles: 增强日志收集器功能和工厂模式支持
- :sparkles: 新增 StudioMax 日志处理器

### 改进
- :zap: 优化 ZipDirectory 函数的并发压缩性能
- :zap: 优化日志收集器的并发处理性能
- :recycle: 重构文件过滤和排序逻辑
- :sparkles: 增强日志处理器的文件匹配和路径处理
- :arrow_up: 清理依赖并移除未使用的模块
- :art: 重构日志格式化器，增强日志输出样式
- :recycle: 调整 service.go 中的目录路径和导入路径
- :recycle: 重命名 hmiserver 为 hmi_server

### 修复
- :bug: 改进 StudioMax 日志处理器的错误处理
- :see_no_evil: 更新 .gitignore，忽略更多测试结果目录

### 测试
- :white_check_mark: 新增 StudioMax 日志处理器的测试用例

### 其他
- :rocket: 添加 GitHub Actions 自动发布工作流
- :rocket: 添加自动发布脚本
- :lock: 更新远程配置链接


## [v0.0.3] - 2025-03-04

### 新特性
- :bookmark: 引入版本管理和构建时信息注入

### 改进
- :sparkles: 优化 HMI 日志处理和文件过滤逻辑
- :recycle: 重命名

### 测试
- :white_check_mark: 新增 HMI 日志处理器测试用例


## [v0.0.2] - 2025-03-04

### 新特性
- ✨ 为更新命令添加 TUI 界面和自动更新支持
- ✨ 为版本命令添加简单模式和 TUI 界面支持
- ✨ 新增程序更新和版本管理功能
- ✨ 新增 Windows 和 Linux 安装脚本
- ✨ 重构上传组件，支持多云存储提供商
- ✨ 重构版本更新逻辑，支持简单模式和 TUI 模式的版本检查
- ✨ 新增版本更新和远程配置管理功能
- ✨ 优化表单交互和日志管理机制
- ✨ 重构项目架构，引入交互式UI和服务层
- :sparkles: 支持 xyz-hmi 日志的收集

### 改进
- :construction_worker: 新增 Makefile 打包目标平台二进制文件的功能
- :recycle: 重命名日志文件保留选项，提高配置语义性
- 🔧 调整日志收集完成后的日志输出逻辑
- ✨ 优化日志收集后处理逻辑
- 🔧 简化命令行上传配置
- 🔧 调整 Cloudreve 分享文件链接参数
- :recycle: 重构命令行模块，拆分功能到独立文件
- :coffin: 移除 UI 相关模块，简化命令行交互逻辑
- 🔧 优化 HMI 服务器日志文件名处理逻辑
- :construction_worker: 优化 Makefile 构建目录结构
- 🔧 增强 Makefile，支持跨平台二进制构建
- 🔧 重构日志收集和处理组件，引入 logrus 日志库
- 🔧 优化 Makefile，增加二进制文件构建选项
- 🔧 支持开发和生产环境的远程配置URL
- 🔧 重构表单和日志组件，优化交互细节
- 🔧 更新依赖并清理无用的第三方库
- 🔧 移除日志配置保存输入组件
- 🔧 精简日志过滤和表单交互逻辑
- 优化日志收集和处理逻辑

### 文档
- 📝 更新 README.md 安装和使用说明
- :memo: 简化 README 中的日志收集命令示例

### 测试
- ✨ 新增日志过滤组件的单元测试和增强处理逻辑

### 其他
- :bookmark: 更新版本为 0.0.2
- 🔧 Update go.mod dependencies
- 初始化提交
