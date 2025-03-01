# logsnap fish completion

function __fish_logsnap_no_subcommand
    set -l cmd (commandline -opc)
    if [ (count $cmd) -eq 1 ]
        return 0
    end
    return 1
end

# 主命令补全
complete -f -c logsnap -n '__fish_logsnap_no_subcommand' -a 'collect' -d '收集指定时间范围内的日志并上传'
complete -f -c logsnap -n '__fish_logsnap_no_subcommand' -a 'update' -d '检查并更新程序到最新版本'
complete -f -c logsnap -n '__fish_logsnap_no_subcommand' -a 'supported-programs' -d '显示支持的程序列表'
complete -f -c logsnap -n '__fish_logsnap_no_subcommand' -a 'version' -d '显示当前版本信息'
complete -f -c logsnap -n '__fish_logsnap_no_subcommand' -a 'completion' -d '生成自动补全脚本'
complete -f -c logsnap -n '__fish_logsnap_no_subcommand' -a 'help' -d '显示帮助信息'

# collect 子命令补全
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'time' -s 't' -d '收集最近多长时间的日志'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'start-time' -s 's' -d '日志收集的开始时间'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'end-time' -s 'e' -d '日志收集的结束时间'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'log-dir' -s 'l' -d '日志目录路径'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'upload' -s 'u' -d '是否上传到云端'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'keep-local-snapshot' -s 'k' -d '是否保留本地日志快照'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'output-dir' -s 'o' -d '输出目录'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'program' -s 'p' -d '要收集的程序日志'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'today' -d '收集今天的日志'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'yesterday' -d '收集昨天的日志'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'this-week' -d '收集本周的日志'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'skip-version-check' -d '跳过版本检查'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'config-dir' -d '配置目录路径'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'simple' -d '使用简单模式，不显示终端动画'
complete -f -c logsnap -n '__fish_seen_subcommand_from collect' -l 'interactive' -s 'I' -d '启用交互模式，通过UI配置选项'

# update 子命令补全
complete -f -c logsnap -n '__fish_seen_subcommand_from update' -l 'force' -s 'f' -d '强制更新，不询问确认'
complete -f -c logsnap -n '__fish_seen_subcommand_from update' -l 'check-only' -s 'c' -d '仅检查是否有更新，不执行更新操作'
complete -f -c logsnap -n '__fish_seen_subcommand_from update' -l 'config-dir' -d '配置目录路径'

# version 子命令补全
complete -f -c logsnap -n '__fish_seen_subcommand_from version' -l 'simple' -s 's' -d '使用简单模式显示版本信息，不使用TUI界面'

# completion 子命令补全
complete -f -c logsnap -n '__fish_seen_subcommand_from completion' -a 'bash' -d '生成 Bash 自动补全脚本'
complete -f -c logsnap -n '__fish_seen_subcommand_from completion' -a 'zsh' -d '生成 Zsh 自动补全脚本'
complete -f -c logsnap -n '__fish_seen_subcommand_from completion' -a 'fish' -d '生成 Fish 自动补全脚本'
complete -f -c logsnap -n '__fish_seen_subcommand_from completion' -a 'powershell' -d '生成 PowerShell 自动补全脚本'
complete -f -c logsnap -n '__fish_seen_subcommand_from completion' -a 'install' -d '自动检测 shell 并安装补全脚本'
