#compdef logsnap

# 将此文件保存为 _logsnap 并放置在 $fpath 中的某个目录下

_logsnap_commands() {
  local -a commands
  commands=(
    'collect:收集指定时间范围内的日志并上传'
    'update:检查并更新程序到最新版本'
    'supported-programs:显示支持的程序列表'
    'version:显示当前版本信息'
    'completion:生成自动补全脚本'
    'help:显示帮助信息'
  )
  _describe -t commands 'logsnap commands' commands
}

_logsnap_collect_options() {
  local -a options
  options=(
    '--time[收集最近多长时间的日志]:时间:(30m 1h 2h 6h 12h 1d 2d 7d)'
    '-t[收集最近多长时间的日志]:时间:(30m 1h 2h 6h 12h 1d 2d 7d)'
    '--start-time[日志收集的开始时间]:开始时间:'
    '-s[日志收集的开始时间]:开始时间:'
    '--end-time[日志收集的结束时间]:结束时间:'
    '-e[日志收集的结束时间]:结束时间:'
    '--log-dir[日志目录路径]:日志目录:_files -/'
    '-l[日志目录路径]:日志目录:_files -/'
    '--upload[是否上传到云端]'
    '-u[是否上传到云端]'
    '--keep-local-snapshot[是否保留本地日志快照]'
    '-k[是否保留本地日志快照]'
    '--output-dir[输出目录]:输出目录:_files -/'
    '-o[输出目录]:输出目录:_files -/'
    '--program[要收集的程序日志]:程序:(xyz-hmi xyz-bin-packing xyz-max-hmi-server xyz-studio-max)'
    '-p[要收集的程序日志]:程序:(xyz-hmi xyz-bin-packing xyz-max-hmi-server xyz-studio-max)'
    '--today[收集今天的日志]'
    '--yesterday[收集昨天的日志]'
    '--this-week[收集本周的日志]'
    '--skip-version-check[跳过版本检查]'
    '--config-dir[配置目录路径]:配置目录:_files -/'
    '--simple[使用简单模式，不显示终端动画]'
    '--interactive[启用交互模式，通过UI配置选项]'
    '-I[启用交互模式，通过UI配置选项]'
  )
  _arguments -s : $options
}

_logsnap_update_options() {
  local -a options
  options=(
    '--force[强制更新，不询问确认]'
    '-f[强制更新，不询问确认]'
    '--check-only[仅检查是否有更新，不执行更新操作]'
    '-c[仅检查是否有更新，不执行更新操作]'
    '--config-dir[配置目录路径]:配置目录:_files -/'
  )
  _arguments -s : $options
}

_logsnap_version_options() {
  local -a options
  options=(
    '--simple[使用简单模式显示版本信息，不使用TUI界面]'
    '-s[使用简单模式显示版本信息，不使用TUI界面]'
  )
  _arguments -s : $options
}

_logsnap_completion_options() {
  local -a options
  options=(
    'bash:生成 Bash 自动补全脚本'
    'zsh:生成 Zsh 自动补全脚本'
    'fish:生成 Fish 自动补全脚本'
    'powershell:生成 PowerShell 自动补全脚本'
    'install:自动检测 shell 并安装补全脚本'
  )
  _describe -t commands 'completion options' options
}

_logsnap() {
  local curcontext="$curcontext" state line
  typeset -A opt_args

  _arguments -C \
    '1: :->command' \
    '*::options:->options'

  case $state in
    command)
      _logsnap_commands
      ;;
    options)
      case $line[1] in
        collect)
          _logsnap_collect_options
          ;;
        update)
          _logsnap_update_options
          ;;
        version)
          _logsnap_version_options
          ;;
        completion)
          _logsnap_completion_options
          ;;
      esac
      ;;
  esac
}

compdef _logsnap logsnap
