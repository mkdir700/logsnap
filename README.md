# LogSnap

LogSnap 是一个高效的轻量级日志收集、打包和上传工具，专为部署人员设计，用于快速收集特定时间范围内的日志文件。

## ✨ 功能特性

- **⏱️ 时间范围收集**：根据指定的时间范围（如最近 30 分钟、1 小时或自定义时间段）收集日志
- **🔍 智能日志解析**：自动识别不同格式的日志文件中的时间戳
- **📦 日志打包**：将收集的日志文件打包成 ZIP 格式，方便传输和存储
- **☁️ 快速分享**：自动将日志文件上传到云端，快速分享给其他同事

## 🧩 架构设计

LogSnap 设计为可扩展的插件式架构，主要由以下组件构成：

1. **Processor**: 日志处理器，负责识别和处理不同类型的日志文件
   - `BaseProcessor`: 基础处理器，提供通用功能
   - 具体处理器实现（如 `BinPackingLogProcessor`）: 处理特定类型的日志

2. **Collector**: 日志收集器，负责按时间范围收集日志

3. **Uploader**: 上传组件，支持将日志上传到不同的存储服务

## 🛠️ 开发指南

### 自定义处理器开发

LogSnap 需要针对特定日志格式开发自定义处理器。以下是开发步骤：

1. **创建新的处理器**：

在 `collector/processor` 目录下新建一个文件夹，例如 `my_custom_processor`，然后创建一个继承自 `BaseProcessor` 的新处理器：

```go
// 创建一个继承自BaseProcessor的新处理器
type MyCustomLogProcessor struct {
    *processor.BaseProcessor
}

func NewMyCustomLogProcessor(logDir string, outputDir string) *MyCustomLogProcessor {
    return &MyCustomLogProcessor{
        BaseProcessor: processor.NewBaseProcessor("我的自定义日志", logDir, outputDir),
    }
}
```

2. **实现文件处理器**：

```go
// 为你的日志格式创建文件处理器
func (p *MyCustomLogProcessor) CreateFileProcessor() []processor.FileProcessorProvider {
    return []processor.FileProcessorProvider{
        NewMyCustomFileProcessorProvider(),
        // 可以添加多个文件处理器
    }
}
```

3. **定义日志文件过滤规则**：

```go
// 创建文件识别和处理规则
type MyCustomFileInfoFilter struct {}

func (f *MyCustomFileInfoFilter) ParseFileInfos(files []string) ([]processor.LogFileInfo, error) {
    // 实现从文件名解析时间戳等信息的逻辑
}

func (f *MyCustomFileInfoFilter) IsMatch(fileName string) bool {
    // 实现识别目标日志文件的逻辑
}
```

4. **定义文件处理器**：

```go
// 创建文件处理器
type MyCustomFileProcessorProvider struct {
    processor.BaseProcessorProvider
}
```

5. **集成到主程序**：

完成处理器后，需要将其注册到主程序的处理器列表中。

参考：https://github.com/mkdir700/logsnap/blob/d603ccf334933e6f6d84442bb14a4a9a4141dd93/collector/factory/factory.go#L45-L61

### 配置文件说明

LogSnap 使用两个主要配置文件，这些文件需要放在云端：

#### 1. 上传配置 (config.json)

```json
{
  "version": "0.0.1",
  "upload_config": {
    "providers": [
      {
        "provider": "s3",
        "endpoint": "https://s3.example.com",
        "bucket": "logs-bucket",
        "region": "us-east-1",
        "folder_path": "xyz-logsnap/snapshots"
      },
      {
        "provider": "webdav",
        "endpoint": "",
        "username": "",
        "password": "",
        "folder_path": "snapshots"
      }
    ],
    "default_provider": "cloudreve"
  }
}
```

> **注意**: 当前版本主要支持 Cloudreve 云盘作为存储和快照分享方案。Cloudreve 作为网盘提供了分享链接失效功能，非常适合临时日志分享需求。虽然配置文件中列出了 S3 和 WebDAV 等其他存储提供商，但目前这些功能尚未完全实现。

#### 2. 下载配置 (download.json)

```json
{
  "latest_versions": {
    "stable": "0.0.2",
    "beta": "0.0.2-beta",
    "dev": "0.0.2-dev"
  },
  "download_urls": {
    "stable": {
      "windows": "https://example.com/logsnap-windows.zip",
      "linux": "https://example.com/logsnap-linux.zip",
      "darwin": "https://example.com/logsnap-darwin.zip"
    }
  },
  "force_update": false,
  "update_message": "发现新版本，建议更新！"
}
```

## 🚀 构建与部署

### 准备工作

1. 修改 `download.json` 和 `config.json` 文件，上传到可访问的HTTP服务器

2. 更新 Makefile 中的配置 URL 变量:

```makefile
UPLOAD_CONFIG_URL ?= "https://your-domain.com/config.json"
DOWNLOAD_CONFIG_URL ?= "https://your-domain.com/download.json"
```

### 构建

```bash
# 常规构建
make build

# 优化构建（更小的二进制文件）
make build-small

# 交叉编译所有支持平台
make build-all

# 打包发布
make package
```

### 修改安装脚本

部署前需要修改 `scripts/install.sh` 中的下载链接：

```bash
# 配置
VERSION_JSON_URL="https://your-domain.com/path/to/download.json"
```

同时确保 `download.json` 中的下载链接指向正确的二进制文件位置。

## 🚀 安装

我们提供了一个通用的安装脚本，可以自动检测您的系统环境并安装 LogSnap。

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

# 收集日志并上传
logsnap c -u
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

## 📝 注意事项

- LogSnap 是一个框架工具，需要根据实际日志格式开发自定义处理器
- 上传功能目前主要支持 Cloudreve 云盘，需要正确配置 Cloudreve 的访问凭证
- Cloudreve 提供了文件分享链接自动失效功能，非常适合临时日志分享需求
- 如需支持其他存储服务，可以通过开发对应的上传组件实现
- 在生产环境使用前，建议先在测试环境验证自定义处理器的功能

## 🤝 贡献

欢迎提交问题报告、功能请求和 Pull Request。对于重大更改，请先打开 issue 讨论您想要更改的内容。

## 📄 许可证

LogSnap 使用 Apache-2.0 许可证。请参阅 [./LICENSE](LICENSE) 文件了解更多信息。
