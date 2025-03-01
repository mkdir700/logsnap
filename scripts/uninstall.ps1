# XYZLogSnap Windows 卸载脚本
# 用法: 
# 在PowerShell中执行: 
# iwr -useb https://your-domain.com/uninstall.ps1 | iex

# 设置错误操作首选项
$ErrorActionPreference = "Stop"

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

# 检查依赖
function Check-Dependencies {
    Write-ColorOutput "检查权限..." "Cyan"
    
    # 检查是否有管理员权限
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if (-not $isAdmin) {
        Write-ColorOutput "未以管理员身份运行。某些操作可能需要管理员权限。" "Yellow"
        $script:NeedAdmin = $true
    }
}

# 查找安装目录
function Find-InstallDirectory {
    Write-ColorOutput "查找安装目录..." "Cyan"
    
    # 可能的安装位置
    $possibleLocations = @(
        "$env:ProgramFiles\XYZLogSnap",
        "$env:LOCALAPPDATA\XYZLogSnap"
    )
    
    foreach ($location in $possibleLocations) {
        if (Test-Path -Path $location) {
            $script:InstallDir = $location
            Write-ColorOutput "找到安装目录: $script:InstallDir" "Green"
            return
        }
    }
    
    Write-ColorOutput "未找到安装目录。XYZLogSnap 可能未安装或已被移除。" "Yellow"
    $script:NotInstalled = $true
}

# 停止运行中的进程
function Stop-RunningProcesses {
    Write-ColorOutput "检查并停止运行中的进程..." "Cyan"
    
    $processes = Get-Process -Name "logsnap" -ErrorAction SilentlyContinue
    
    if ($processes) {
        Write-ColorOutput "发现正在运行的 XYZLogSnap 进程，正在停止..." "Yellow"
        
        foreach ($process in $processes) {
            try {
                $process | Stop-Process -Force
                Write-ColorOutput "已停止进程 ID: $($process.Id)" "Green"
            } catch {
                Write-ColorOutput "无法停止进程 ID: $($process.Id). $($_.Exception.Message)" "Red"
            }
        }
    } else {
        Write-ColorOutput "未发现正在运行的 XYZLogSnap 进程。" "Green"
    }
}

# 从环境变量中移除
function Remove-FromPath {
    Write-ColorOutput "从环境变量中移除..." "Cyan"
    
    if ($script:NotInstalled) {
        return
    }
    
    # 从 PATH 中移除
    $userPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
    $machinePath = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
    
    # 用户级别 PATH
    if ($userPath -like "*$script:InstallDir*") {
        $newUserPath = ($userPath -split ';' | Where-Object { $_ -ne $script:InstallDir -and $_ -ne "$script:InstallDir\" }) -join ';'
        [System.Environment]::SetEnvironmentVariable("PATH", $newUserPath, "User")
        Write-ColorOutput "已从用户 PATH 中移除。" "Green"
    }
    
    # 系统级别 PATH (需要管理员权限)
    if (-not $script:NeedAdmin -and $machinePath -like "*$script:InstallDir*") {
        $newMachinePath = ($machinePath -split ';' | Where-Object { $_ -ne $script:InstallDir -and $_ -ne "$script:InstallDir\" }) -join ';'
        [System.Environment]::SetEnvironmentVariable("PATH", $newMachinePath, "Machine")
        Write-ColorOutput "已从系统 PATH 中移除。" "Green"
    }
}

# 删除桌面快捷方式
function Remove-Shortcut {
    Write-ColorOutput "删除桌面快捷方式..." "Cyan"
    
    $desktopPath = [System.Environment]::GetFolderPath("Desktop")
    $shortcutPath = "$desktopPath\XYZLogSnap.lnk"
    
    if (Test-Path -Path $shortcutPath) {
        Remove-Item -Path $shortcutPath -Force
        Write-ColorOutput "已删除桌面快捷方式。" "Green"
    } else {
        Write-ColorOutput "未找到桌面快捷方式。" "Yellow"
    }
}

# 删除安装目录
function Remove-InstallDirectory {
    Write-ColorOutput "删除安装目录..." "Cyan"
    
    if ($script:NotInstalled) {
        return
    }
    
    try {
        Remove-Item -Path $script:InstallDir -Recurse -Force
        Write-ColorOutput "已删除安装目录: $script:InstallDir" "Green"
    } catch {
        Write-ColorOutput "无法删除安装目录: $($_.Exception.Message)" "Red"
        Write-ColorOutput "请手动删除目录: $script:InstallDir" "Yellow"
    }
}

# 清理注册表
function Clean-Registry {
    Write-ColorOutput "清理注册表..." "Cyan"
    
    # 可能的注册表路径
    $registryPaths = @(
        "HKCU:\Software\XYZLogSnap",
        "HKLM:\Software\XYZLogSnap"
    )
    
    foreach ($path in $registryPaths) {
        if (Test-Path -Path $path) {
            try {
                Remove-Item -Path $path -Recurse -Force
                Write-ColorOutput "已删除注册表项: $path" "Green"
            } catch {
                Write-ColorOutput "无法删除注册表项 $path: $($_.Exception.Message)" "Red"
            }
        }
    }
}

# 主函数
function Main {
    Write-ColorOutput "=== XYZLogSnap Windows 卸载程序 ===" "Green"
    
    $script:NeedAdmin = $false
    $script:InstallDir = ""
    $script:NotInstalled = $false
    
    Check-Dependencies
    Find-InstallDirectory
    
    if ($script:NotInstalled) {
        Write-ColorOutput "XYZLogSnap 似乎未安装，卸载过程将终止。" "Yellow"
        return
    }
    
    # 确认卸载
    Write-ColorOutput "您确定要卸载 XYZLogSnap 吗? 这将删除所有程序文件和设置。" "Yellow"
    Write-ColorOutput "输入 'Y' 继续卸载，或任意其他键取消: " "Yellow" -NoNewline
    
    $confirmation = Read-Host
    if ($confirmation -ne "Y" -and $confirmation -ne "y") {
        Write-ColorOutput "卸载已取消。" "Cyan"
        return
    }
    
    Stop-RunningProcesses
    Remove-FromPath
    Remove-Shortcut
    Remove-InstallDirectory
    Clean-Registry
    
    Write-ColorOutput "XYZLogSnap 已成功卸载!" "Green"
}

# 执行主函数
Main 