#!/bin/bash

# 检测当前使用的 shell
detect_shell() {
    # 首先尝试从环境变量获取
    shell="$SHELL"
    
    # 在 Windows 上，尝试检测 PowerShell
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OSTYPE" == "cygwin" || -z "$shell" ]]; then
        # 检查是否在 PowerShell 中运行
        if [[ -n "$PSModulePath" ]]; then
            echo "powershell"
            return
        fi
        # 默认返回 bash，因为它是最常见的
        echo "bash"
        return
    fi
    
    # 提取 shell 名称
    shell=$(basename "$shell")
    echo "$shell"
}

# 安装 Bash 补全脚本
install_bash_completion() {
    echo "安装 Bash 补全脚本..."
    
    home_dir="$HOME"
    completion_file="$home_dir/.logsnap-completion.bash"
    
    # 复制补全脚本
    cp "$(dirname "$0")/bash-completion.sh" "$completion_file"
    chmod 644 "$completion_file"
    
    # 检查 .bashrc 中是否已有加载命令
    bashrc_path="$home_dir/.bashrc"
    
    # 要添加的加载命令
    load_cmd="\n# logsnap 自动补全\n[ -f $completion_file ] && source $completion_file\n"
    
    # 如果 .bashrc 不存在或不包含加载命令，则添加
    if [[ ! -f "$bashrc_path" ]] || ! grep -q "logsnap-completion.bash" "$bashrc_path"; then
        echo -e "$load_cmd" >> "$bashrc_path"
        echo "已将补全脚本加载命令添加到 ~/.bashrc"
        echo "请运行 'source ~/.bashrc' 或重新打开终端以启用补全功能"
    else
        echo "补全脚本已配置在 ~/.bashrc 中"
    fi
}

# 安装 Zsh 补全脚本
install_zsh_completion() {
    echo "安装 Zsh 补全脚本..."
    
    home_dir="$HOME"
    
    # 确保 zsh 补全目录存在
    zsh_completion_dir="$home_dir/.zsh/completion"
    mkdir -p "$zsh_completion_dir"
    
    # 创建补全脚本文件
    completion_file="$zsh_completion_dir/_logsnap"
    
    # 复制补全脚本
    cp "$(dirname "$0")/zsh-completion.sh" "$completion_file"
    chmod 644 "$completion_file"
    
    # 检查 .zshrc 中是否已有加载命令
    zshrc_path="$home_dir/.zshrc"
    
    # 要添加的加载命令
    load_cmd="\n# logsnap 自动补全\nfpath=($zsh_completion_dir \$fpath)\nautoload -Uz compinit && compinit\n"
    
    # 如果 .zshrc 不存在或不包含加载命令，则添加
    if [[ ! -f "$zshrc_path" ]] || ! grep -q "$zsh_completion_dir" "$zshrc_path"; then
        echo -e "$load_cmd" >> "$zshrc_path"
        echo "已将补全脚本加载命令添加到 ~/.zshrc"
        echo "请运行 'source ~/.zshrc' 或重新打开终端以启用补全功能"
    else
        echo "补全脚本已配置在 ~/.zshrc 中"
    fi
}

# 安装 Fish 补全脚本
install_fish_completion() {
    echo "安装 Fish 补全脚本..."
    
    home_dir="$HOME"
    
    # 确保 fish 补全目录存在
    fish_completion_dir="$home_dir/.config/fish/completions"
    mkdir -p "$fish_completion_dir"
    
    # 创建补全脚本文件
    completion_file="$fish_completion_dir/logsnap.fish"
    
    # 复制补全脚本
    cp "$(dirname "$0")/fish-completion.fish" "$completion_file"
    chmod 644 "$completion_file"
    
    echo "补全脚本已安装到 ~/.config/fish/completions/logsnap.fish"
    echo "Fish 会自动加载此目录中的补全脚本，无需额外配置"
}

# 安装 PowerShell 补全脚本
install_powershell_completion() {
    echo "安装 PowerShell 补全脚本..."
    
    # 获取 PowerShell 配置文件路径
    profile_path=$(powershell -Command "echo \$PROFILE.CurrentUserAllHosts" 2>/dev/null)
    
    if [[ -z "$profile_path" ]]; then
        echo "错误：无法获取 PowerShell 配置文件路径"
        return 1
    fi
    
    # 确保配置文件目录存在
    profile_dir=$(dirname "$profile_path")
    mkdir -p "$profile_dir"
    
    # 检查配置文件中是否已有补全脚本
    if [[ -f "$profile_path" ]] && grep -q "logsnap PowerShell completion" "$profile_path"; then
        echo "补全脚本已配置在 PowerShell 配置文件中"
        return 0
    fi
    
    # 复制补全脚本内容到配置文件
    cat "$(dirname "$0")/powershell-completion.ps1" >> "$profile_path"
    
    echo "已将补全脚本添加到 PowerShell 配置文件: $profile_path"
    echo "请重新打开 PowerShell 或运行 '. \$PROFILE.CurrentUserAllHosts' 以启用补全功能"
}

# 主函数
main() {
    shell=$(detect_shell)
    echo "检测到当前 shell 为: $shell"
    
    case "$shell" in
        bash)
            install_bash_completion
            ;;
        zsh)
            install_zsh_completion
            ;;
        fish)
            install_fish_completion
            ;;
        powershell|pwsh)
            install_powershell_completion
            ;;
        *)
            echo "不支持的 shell 类型: $shell"
            echo "支持的 shell 类型: bash, zsh, fish, powershell"
            exit 1
            ;;
    esac
    
    echo "补全脚本安装成功！"
}

main "$@"
