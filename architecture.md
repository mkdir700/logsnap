
## 1. 整体架构

  

```mermaid

graph TD

%% 主要组件

CLI[命令行界面CLI]

Service[核心服务Service]

Collector[日志收集器Collector]

Processor[日志处理器Processor]

Uploader[上传服务Uploader]

%% 连接关系

CLI -->|解析参数| Service

Service -->|创建收集器| Collector

Service -->|配置上传| Uploader

Collector -->|使用| Processor

%% 样式

classDef core fill:#f9f,stroke:#333,stroke-width:2px;

class CLI,Service,Collector,Processor,Uploader core;

```

  

## 2. 核心模块说明

  

### 2.1 命令行界面 (CLI)

  

```mermaid

graph LR

CLI[命令行界面CLI] --> Collect[收集命令collect]

CLI --> Update[更新命令update]

CLI --> Version[版本命令version]

CLI --> Support[支持命令supported-programs]

CLI --> Completion[补全命令completion]

Collect --> TimeOptions[时间选项]

Collect --> UploadOptions[上传选项]

Collect --> OutputOptions[输出选项]

TimeOptions --> RelativeTime[相对时间30m, 1h, 2d]

TimeOptions --> AbsoluteTime[绝对时间开始/结束时间]

TimeOptions --> Shortcuts[快捷选项今天/昨天/本周]

classDef command fill:#bbf,stroke:#333,stroke-width:1px;

classDef options fill:#ddf,stroke:#333,stroke-width:1px;

class CLI command;

class Collect,Update,Version,Support,Completion command;

class TimeOptions,UploadOptions,OutputOptions,RelativeTime,AbsoluteTime,Shortcuts options;

```

  

### 2.2 核心服务 (Service)

  

```mermaid

graph TD

Service[核心服务Service] --> Config[配置管理]

Service --> UploadManager[上传管理器]

Service --> CollectService[收集服务]

Service --> VersionCheck[版本检查]

Config --> LocalConfig[本地配置]

Config --> RemoteConfig[远程配置]

UploadManager --> UploadRequest[上传请求]

UploadManager --> ProgressReport[进度报告]

CollectService --> ProcessorManagement[处理器管理]

CollectService --> TimeRangeProcess[时间范围处理]

classDef service fill:#fbb,stroke:#333,stroke-width:1px;

classDef component fill:#fdd,stroke:#333,stroke-width:1px;

class Service service;

class Config,UploadManager,CollectService,VersionCheck service;

class LocalConfig,RemoteConfig,UploadRequest,ProgressReport,ProcessorManagement,TimeRangeProcess component;

```

  

### 2.3 日志收集器 (Collector)

  

```mermaid

graph TD

Collector[日志收集器Collector] --> Processors[处理器列表]

Collector --> ZipUtils[ZIP工具]

Collector --> ParallelProcess[并行处理]

Processors --> ProcessorInterface[处理器接口LogProcessor]

ProcessorInterface --> GetName[获取名称]

ProcessorInterface --> GetLogPath[获取日志路径]

ProcessorInterface --> Collect[收集日志]

ParallelProcess --> WaitGroup[等待组]

ParallelProcess --> ResultChannel[结果通道]

classDef collector fill:#bfb,stroke:#333,stroke-width:1px;

classDef interface fill:#dfd,stroke:#333,stroke-width:1px;

classDef method fill:#efe,stroke:#333,stroke-width:1px;

class Collector collector;

class Processors,ZipUtils,ParallelProcess collector;

class ProcessorInterface interface;

class GetName,GetLogPath,Collect,WaitGroup,ResultChannel method;

```

  

### 2.4 日志处理器 (Processor)

  

```mermaid

graph TD

Factory[处理器工厂] --> CreateProcessor[创建处理器]

Factory --> GetSupportedTypes[获取支持类型]

CreateProcessor --> HMI[HMI日志处理器]

CreateProcessor --> HMIServer[HMI服务器日志处理器]

CreateProcessor --> StudioMax[StudioMax日志处理器]

CreateProcessor --> BinPacking[BinPacking日志处理器]

CreateProcessor --> VisionLog[VisionLogViewer日志处理器]

CreateProcessor --> RobotDriver[RobotDriver日志处理器]

classDef factory fill:#bbf,stroke:#333,stroke-width:1px;

classDef processor fill:#ddf,stroke:#333,stroke-width:1px;

class Factory factory;

class CreateProcessor,GetSupportedTypes factory;

class HMI,HMIServer,StudioMax,BinPacking,VisionLog,RobotDriver processor;

```

  

### 2.5 上传服务 (Uploader)

  

```mermaid

graph TD

Uploader[上传服务Uploader] --> UploadInterface[云上传接口CloudUploaderInterface]

Uploader --> ProviderSelection[提供商选择]

UploadInterface --> Upload[上传方法]

ProviderSelection --> S3[S3上传器]

ProviderSelection --> WebDAV[WebDAV上传器]

ProviderSelection --> Local[本地存储上传器]

ProviderSelection --> Cloudreve[Cloudreve上传器]

classDef uploader fill:#fbf,stroke:#333,stroke-width:1px;

classDef interface fill:#fdf,stroke:#333,stroke-width:1px;

classDef provider fill:#fef,stroke:#333,stroke-width:1px;

class Uploader uploader;

class UploadInterface,ProviderSelection interface;

class S3,WebDAV,Local,Cloudreve provider;

```

  

## 3. 数据流

  

```mermaid

graph LR

Input[用户输入参数] -->|解析| TimeRange[解析时间范围]

TimeRange -->|创建| ProcessorInst[创建处理器实例]

ProcessorInst -->|处理| CollectFiles[收集日志文件]

CollectFiles -->|过滤| FilterTime[过滤时间范围]

FilterTime -->|压缩| ZipFiles[打包为ZIP]

ZipFiles -->|上传| UploadCloud[上传到云存储]

UploadCloud -->|返回| ReturnURL[返回URL链接]

classDef flow fill:#fff,stroke:#333,stroke-width:1px;

class Input,TimeRange,ProcessorInst,CollectFiles,FilterTime,ZipFiles,UploadCloud,ReturnURL flow;

```

  

## 4. 主要功能模块

  

### 4.1 时间范围收集

  

```mermaid

graph TD

TimeCollection[时间范围收集] --> RelativeTime[相对时间]

TimeCollection --> AbsoluteTime[绝对时间范围]

TimeCollection --> Shortcuts[便捷选项]

RelativeTime --> Minutes[分钟30m]

RelativeTime --> Hours[小时1h]

RelativeTime --> Days[天2d]

AbsoluteTime --> StartTime[开始时间]

AbsoluteTime --> EndTime[结束时间]

Shortcuts --> Today[今天]

Shortcuts --> Yesterday[昨天]

Shortcuts --> ThisWeek[本周]

classDef time fill:#ffd,stroke:#333,stroke-width:1px;

classDef option fill:#ffe,stroke:#333,stroke-width:1px;

class TimeCollection time;

class RelativeTime,AbsoluteTime,Shortcuts time;

class Minutes,Hours,Days,StartTime,EndTime,Today,Yesterday,ThisWeek option;

```

  

### 4.2 智能日志解析

  

```mermaid

graph TD

LogParsing[智能日志解析] --> FormatDetection[格式检测]

LogParsing --> TimeFiltering[时间过滤]

LogParsing --> ContentExtraction[内容提取]

FormatDetection --> HMIFormat[HMI格式]

FormatDetection --> ServerFormat[服务器格式]

FormatDetection --> StudioFormat[Studio格式]

TimeFiltering --> ParseTimestamp[解析时间戳]

TimeFiltering --> CompareRange[比较时间范围]

ContentExtraction --> LineProcessing[行处理]

ContentExtraction --> StructurePreserve[结构保留]

classDef parsing fill:#dff,stroke:#333,stroke-width:1px;

classDef feature fill:#eff,stroke:#333,stroke-width:1px;

class LogParsing parsing;

class FormatDetection,TimeFiltering,ContentExtraction parsing;

class HMIFormat,ServerFormat,StudioFormat,ParseTimestamp,CompareRange,LineProcessing,StructurePreserve feature;

```

  

### 4.3 日志打包

  

```mermaid

graph TD

LogPackaging[日志打包] --> ZipCompression[ZIP压缩]

LogPackaging --> DirectoryStructure[目录结构保留]

LogPackaging --> ParallelProcessing[并行处理]

ZipCompression --> CreateArchive[创建归档]

ZipCompression --> AddFiles[添加文件]

ZipCompression --> Verification[验证完整性]

DirectoryStructure --> PathMapping[路径映射]

DirectoryStructure --> RelativePaths[相对路径]

ParallelProcessing --> Goroutines[Go协程]

ParallelProcessing --> WaitGroups[等待组]

ParallelProcessing --> Channels[通道]

classDef packaging fill:#fdb,stroke:#333,stroke-width:1px;

classDef feature fill:#fed,stroke:#333,stroke-width:1px;

class LogPackaging packaging;

class ZipCompression,DirectoryStructure,ParallelProcessing packaging;

class CreateArchive,AddFiles,Verification,PathMapping,RelativePaths,Goroutines,WaitGroups,Channels feature;

```

  

### 4.4 快速分享

  

```mermaid

graph TD

Sharing[快速分享] --> UploadMethods[上传方式]

Sharing --> URLGeneration[URL生成]

Sharing --> ProgressReporting[进度报告]

UploadMethods --> S3Upload[S3上传]

UploadMethods --> WebDAVUpload[WebDAV上传]

UploadMethods --> LocalUpload[本地存储]

UploadMethods --> CloudreveUpload[Cloudreve上传]

URLGeneration --> SignedURL[签名URL]

URLGeneration --> ExpirationTime[过期时间]

ProgressReporting --> Percentage[百分比]

ProgressReporting --> SpeedCalculation[速度计算]

classDef sharing fill:#dbf,stroke:#333,stroke-width:1px;

classDef feature fill:#edf,stroke:#333,stroke-width:1px;

class Sharing sharing;

class UploadMethods,URLGeneration,ProgressReporting sharing;

class S3Upload,WebDAVUpload,LocalUpload,CloudreveUpload,SignedURL,ExpirationTime,Percentage,SpeedCalculation feature;

```

  

## 5. 扩展性设计

  

### 5.1 处理器扩展

  

```mermaid

graph TD

ProcessorExtension[处理器扩展] --> DefineType[定义处理器类型]

ProcessorExtension --> ImplementInterface[实现LogProcessor接口]

ProcessorExtension --> RegisterFactory[注册处理器工厂]

ImplementInterface --> GetNameMethod[GetName方法]

ImplementInterface --> GetLogPathMethod[GetLogPath方法]

ImplementInterface --> CollectMethod[Collect方法]

RegisterFactory --> UpdateRegistry[更新注册表]

RegisterFactory --> CreateFactoryImpl[创建工厂实现]

classDef extension fill:#bdf,stroke:#333,stroke-width:1px;

classDef step fill:#def,stroke:#333,stroke-width:1px;

class ProcessorExtension extension;

class DefineType,ImplementInterface,RegisterFactory extension;

class GetNameMethod,GetLogPathMethod,CollectMethod,UpdateRegistry,CreateFactoryImpl step;

```

  

### 5.2 上传方式扩展

  

```mermaid

graph TD

UploaderExtension[上传方式扩展] --> ImplementInterface[实现CloudUploaderInterface]

UploaderExtension --> AddProvider[添加提供商支持]

ImplementInterface --> UploadMethod[Upload方法]

AddProvider --> DefineProvider[定义提供商类型]

AddProvider --> UpdateUploaderSwitch[更新上传器选择逻辑]

classDef extension fill:#fbd,stroke:#333,stroke-width:1px;

classDef step fill:#fde,stroke:#333,stroke-width:1px;

class UploaderExtension extension;

class ImplementInterface,AddProvider extension;

class UploadMethod,DefineProvider,UpdateUploaderSwitch step;

```

  

## 6. 配置管理

  

```mermaid

graph TD

ConfigManagement[配置管理] --> LocalConfig[本地配置]

ConfigManagement --> RemoteConfig[远程配置]

ConfigManagement --> VersionManagement[版本管理]

LocalConfig --> ConfigFile[配置文件]

LocalConfig --> ConfigDir[配置目录]

RemoteConfig --> RemoteURL[远程URL]

RemoteConfig --> UpdateInterval[更新间隔]

VersionManagement --> VersionCheck[版本检查]

VersionManagement --> AutoUpdate[自动更新]

classDef config fill:#ddf,stroke:#333,stroke-width:1px;

classDef feature fill:#eef,stroke:#333,stroke-width:1px;

class ConfigManagement config;

class LocalConfig,RemoteConfig,VersionManagement config;

class ConfigFile,ConfigDir,RemoteURL,UpdateInterval,VersionCheck,AutoUpdate feature;

```