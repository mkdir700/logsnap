# XYZLogSnap Windows 安装脚本
# 用法: 
# 在PowerShell中执行: 
# iwr -useb https://raw.githubusercontent.com/mkdir700/xyz-logsnap-release/master/scripts/install.ps1 | iex

# 设置错误操作首选项
$ErrorActionPreference = "Stop"

# 配置
$VersionJsonUrl = "https://example.com/path/to/version.json"

# 颜色定义
function Write-ColorOutput {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Message,
        
        [Parameter(Mandatory = $false)]
        [string]$ForegroundColor = "White"
    )
    
    $previousColor = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    Write-Output $Message
    $host.UI.RawUI.ForegroundColor = $previousColor
}

# 检测系统架构
function Detect-Architecture {
    Write-ColorOutput "检测系统架构..." "Cyan"
    
    $arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")
    
    if ($arch -eq "AMD64") {
        $script:Arch = "amd64"
    } elseif ($arch -eq "ARM64") {
        $script:Arch = "arm64"
    } elseif ($arch -eq "X86") {
        $script:Arch = "386"
    } else {
        Write-ColorOutput "不支持的架构: $arch" "Red"
        exit 1
    }
    
    Write-ColorOutput "检测到架构: $script:Arch" "Green"
    $script:OS = "windows"
}

# 获取稳定版本信息
function Get-StableVersion {
    Write-ColorOutput "获取稳定版本信息..." "Cyan"
    
    # 调用函数获取操作系统和架构信息
    Detect-Architecture
    
    try {
        # 从配置文件获取稳定版本信息
        Write-ColorOutput "获取版本信息..." "Cyan"
        $versionJson = Invoke-RestMethod -Uri $VersionJsonUrl -ErrorAction Stop
        
        # 调试输出
        # Write-ColorOutput "调试信息: 获取到的配置文件内容:" "Yellow"
        # $versionJson | ConvertTo-Json -Depth 10 | Write-Output
        
        # 提取稳定版本号
        $script:Version = $versionJson.latest_versions.stable
        
        if ([string]::IsNullOrEmpty($script:Version)) {
            Write-ColorOutput "无法从配置文件中提取稳定版本号，安装失败" "Red"
            Write-ColorOutput "尝试手动提取版本号..." "Yellow"
            $versionJson.latest_versions | ConvertTo-Json | Write-Output
            exit 1
        }
        
        Write-ColorOutput "获取到稳定版本: $script:Version" "Green"
        
        # 根据操作系统获取对应的下载链接
        $script:DownloadUrl = $versionJson.download_urls.stable.windows
        
        if ([string]::IsNullOrEmpty($script:DownloadUrl)) {
            Write-ColorOutput "无法获取 Windows 系统的下载链接，安装失败" "Red"
            Write-ColorOutput "调试信息: 下载链接部分:" "Yellow"
            $versionJson.download_urls.stable | ConvertTo-Json | Write-Output
            exit 1
        }
        
        Write-ColorOutput "下载链接: $script:DownloadUrl" "Cyan"
    }
    catch {
        Write-ColorOutput "获取版本信息失败: $($_.Exception.Message)" "Red"
        Write-ColorOutput "请检查网络连接或配置文件格式" "Red"
        exit 1
    }
}

# 检查依赖
function Check-Dependencies {
    Write-ColorOutput "检查依赖..." "Cyan"
    
    # 检查是否有管理员权限
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if (-not $isAdmin) {
        Write-ColorOutput "未以管理员身份运行。某些操作可能需要管理员权限。" "Yellow"
        $script:NeedAdmin = $true
    }
    
    # 检查 PowerShell 版本
    if ($PSVersionTable.PSVersion.Major -lt 5) {
        Write-ColorOutput "需要 PowerShell 5.0 或更高版本。" "Red"
        exit 1
    }
    
    Write-ColorOutput "所有依赖已满足。" "Green"
}

# 创建安装目录
function Create-InstallDirectory {
    Write-ColorOutput "创建安装目录..." "Cyan"
    
    # 设置安装目录
    $script:InstallDir = "$env:LOCALAPPDATA\XYZLogSnap"
    
    # 如果有管理员权限，则安装到程序文件目录
    if (-not $script:NeedAdmin) {
        $script:InstallDir = "$env:ProgramFiles\XYZLogSnap"
    }
    
    # 创建目录
    if (-not (Test-Path -Path $script:InstallDir)) {
        New-Item -ItemType Directory -Path $script:InstallDir -Force | Out-Null
    }
    
    Write-ColorOutput "安装目录: $script:InstallDir" "Green"
}

# 下载并安装
function Download-AndInstall {
    Write-ColorOutput "下载 XYZLogSnap v$script:Version..." "Cyan"
    
    # 设置临时目录
    $tempDir = [System.IO.Path]::GetTempPath()
    $tempFolder = Join-Path -Path $tempDir -ChildPath ([System.Guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $tempFolder -Force | Out-Null
    
    # 设置下载文件名
    $outputFile = Join-Path -Path $tempFolder -ChildPath "logsnap.zip"
    
    Write-ColorOutput "从 $script:DownloadUrl 下载中..." "Cyan"
    
    # 下载文件，最多重试3次
    $maxRetries = 3
    $retryCount = 0
    $downloadSuccess = $false
    
    while ($retryCount -lt $maxRetries -and -not $downloadSuccess) {
        try {
            # 下载文件
            [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
            Invoke-WebRequest -Uri $script:DownloadUrl -OutFile $outputFile -UseBasicParsing
            $downloadSuccess = $true
        }
        catch {
            $retryCount++
            if ($retryCount -lt $maxRetries) {
                Write-ColorOutput "下载失败，正在重试 ($retryCount/$maxRetries)..." "Yellow"
                Start-Sleep -Seconds 2
            } else {
                Write-ColorOutput "下载失败。请检查网络连接或版本是否存在。" "Red"
                Write-ColorOutput "下载URL: $script:DownloadUrl" "Red"
                Write-ColorOutput "错误信息: $($_.Exception.Message)" "Red"
                exit 1
            }
        }
    }
    
    # 验证下载文件
    if (-not (Test-Path -Path $outputFile)) {
        Write-ColorOutput "下载文件不存在: $outputFile" "Red"
        exit 1
    }
    
    # 验证文件是否为ZIP格式
    try {
        $fileBytes = [System.IO.File]::ReadAllBytes($outputFile)
        $zipSignature = [byte[]]@(80, 75, 3, 4) # ZIP文件头部标识 (PK\003\004)
        $isZipFile = $true
        
        for ($i = 0; $i -lt [Math]::Min($zipSignature.Length, $fileBytes.Length); $i++) {
            if ($fileBytes[$i] -ne $zipSignature[$i]) {
                $isZipFile = $false
                break
            }
        }
        
        if (-not $isZipFile) {
            Write-ColorOutput "错误：下载的文件不是有效的ZIP格式" "Red"
            exit 1
        }
    }
    catch {
        Write-ColorOutput "无法验证文件格式: $($_.Exception.Message)" "Red"
        exit 1
    }
    
    # 解压文件
    Write-ColorOutput "解压文件..." "Cyan"
    try {
        Expand-Archive -Path $outputFile -DestinationPath $script:InstallDir -Force
    }
    catch {
        Write-ColorOutput "解压失败。文件可能已损坏或格式不正确。" "Red"
        Write-ColorOutput "错误信息: $($_.Exception.Message)" "Red"
        exit 1
    }
    
    # 清理临时文件
    try {
        Remove-Item -Path $tempFolder -Recurse -Force -ErrorAction SilentlyContinue
    }
    catch {
        Write-ColorOutput "清理临时文件失败，但安装可能已成功。" "Yellow"
    }
    
    Write-ColorOutput "文件下载和解压完成。" "Green"
}

# 添加到环境变量
function Add-ToPath {
    Write-ColorOutput "添加到环境变量..." "Cyan"
    
    $exePath = "$script:InstallDir\logsnap.exe"
    
    # 检查可执行文件是否存在
    if (-not (Test-Path -Path $exePath)) {
        Write-ColorOutput "未找到可执行文件 $exePath" "Red"
        exit 1
    }
    
    # 添加到 PATH
    $userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    $machinePath = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
    
    if ($script:NeedAdmin) {
        # 用户级别 PATH
        if ($userPath -notlike "*$script:InstallDir*") {
            [System.Environment]::SetEnvironmentVariable("PATH", "$userPath;$script:InstallDir", "User")
            $env:PATH = "$env:PATH;$script:InstallDir"
        }
    } else {
        # 系统级别 PATH
        if ($machinePath -notlike "*$script:InstallDir*") {
            [System.Environment]::SetEnvironmentVariable("PATH", "$machinePath;$script:InstallDir", "Machine")
            $env:PATH = "$env:PATH;$script:InstallDir"
        }
    }
    
    Write-ColorOutput "已添加到环境变量。" "Green"
}

# 创建桌面快捷方式
function Create-Shortcut {
    Write-ColorOutput "创建桌面快捷方式..." "Cyan"
    
    $desktopPath = [System.Environment]::GetFolderPath("Desktop")
    $shortcutPath = "$desktopPath\XYZLogSnap.lnk"
    $exePath = "$script:InstallDir\logsnap.exe"
    
    $WshShell = New-Object -ComObject WScript.Shell
    $Shortcut = $WshShell.CreateShortcut($shortcutPath)
    $Shortcut.TargetPath = $exePath
    $Shortcut.Description = "XYZLogSnap - 高效的日志收集工具"
    $Shortcut.WorkingDirectory = $script:InstallDir
    $Shortcut.Save()
    
    Write-ColorOutput "桌面快捷方式已创建。" "Green"
}

# 验证安装
function Verify-Installation {
    Write-ColorOutput "验证安装..." "Cyan"
    
    $exePath = "$script:InstallDir\logsnap.exe"
    
    if (Test-Path -Path $exePath) {
        Write-ColorOutput "XYZLogSnap 稳定版 v$script:Version 已成功安装!" "Green"
        
        # 尝试获取版本信息
        try {
            $versionInfo = & $exePath version 2>&1
            Write-ColorOutput "版本信息: $versionInfo" "Green"
        } catch {
            Write-ColorOutput "无法获取版本信息: $($_.Exception.Message)" "Yellow"
        }
        
        Write-ColorOutput "`n使用方法示例:" "Cyan"
        Write-ColorOutput "  logsnap collect - 收集最近30分钟的日志" "Yellow"
        Write-ColorOutput "  logsnap collect -u - 收集最近30分钟的日志并上传云端" "Yellow"
        Write-ColorOutput "  logsnap collect --time 1h - 收集最近1小时的日志" "Yellow"
        Write-ColorOutput "  logsnap collect --start-time `"2023-03-01 10:00:00`" --end-time `"2023-03-01 11:00:00`" - 收集指定时间范围的日志" "Yellow"
    } else {
        Write-ColorOutput "安装似乎失败，请手动检查。" "Red"
        Write-ColorOutput "尝试手动运行: $exePath version" "Yellow"
        exit 1
    }
}

# 主函数
function Main {
    Write-ColorOutput "=== XYZLogSnap 稳定版安装程序 ===" "Green"
    
    $script:NeedAdmin = $false
    $script:Arch = ""
    $script:OS = ""
    $script:InstallDir = ""
    $script:Version = ""
    $script:DownloadUrl = ""
    
    Check-Dependencies
    Get-StableVersion
    Create-InstallDirectory
    Download-AndInstall
    Add-ToPath
    Create-Shortcut
    Verify-Installation
    
    Write-ColorOutput "安装过程完成!" "Green"
}

# 执行主函数
Main