# XYZLogSnap

XYZLogSnap 是一个高效的日志收集、打包和上传工具，专为部署人员设计，用于快速收集特定时间范围内的日志文件。

## ✨ 功能特性

- **⏱️ 时间范围收集**：根据指定的时间范围（如最近 30 分钟、1 小时或自定义时间段）收集日志
- **🔍 智能日志解析**：自动识别不同格式的日志文件中的时间戳
- **📦 日志打包**：将收集的日志文件打包成 ZIP 格式，方便传输和存储
- **☁️ 快速分享**：自动将日志文件上传到云端，快速分享给其他同事

## 🚀 安装

我们提供了一个通用的安装脚本，可以自动检测您的系统环境并安装 XYZLogSnap。

### 🐧 Linux 安装

```bash
curl -sSL "https://example.com/install.sh?=$(date +%s)" | bash
```

### 🪟 Windows 安装

在 Windows 系统上，我们提供了 PowerShell 安装脚本：

```powershell
iwr -useb "https://example.com/install.ps1?=$(Get-Random)" | iex
```

> 💡 提示：右键点击 PowerShell 图标，选择"以管理员身份运行"，然后执行上述命令。

## 📖 使用方法

### 🔰 基本用法

```bash
# 收集最近30分钟的日志
logsnap c

# 收集最近1小时的日志
logsnap c --time 1h

# 收集指定时间范围的日志
logsnap c --start-time "2023-03-01 10:00:00" --end-time "2023-03-01 11:00:00"

# 收集日志但不上传
logsnap c --upload=false
```

### 🎮 命令行选项

- `--time, -t`：指定收集最近多长时间的日志（例如：30m, 1h, 2d）
- `--start-time, -s`：日志收集的开始时间（格式：YYYY-MM-DD HH:MM:SS）
- `--end-time, -e`：日志收集的结束时间（格式：YYYY-MM-DD HH:MM:SS，默认为当前时间）
- `--upload, -u`：是否上传收集的日志（默认：false）
- `--keep-local-snapshot, -k`：是否保留本地日志快照（默认：false）

## 🗑️ 卸载

如果您需要卸载 XYZLogSnap，可以执行以下命令：

```bash
# 删除可执行文件
sudo rm -f /usr/local/bin/logsnap

# 删除安装目录
rm -rf ~/.logsnap
```
