#! /bin/bash

_logsnap_completion() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  
  # 完成 logsnap 命令的补全
  if [[ ${COMP_CWORD} -eq 1 ]]; then
    opts="collect update version supported-programs completion help"
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi

  # 根据子命令提供不同的补全选项
  case "${COMP_WORDS[1]}" in
    collect)
      # 处理特定参数的补全
      if [[ ${prev} == "--time" || ${prev} == "-t" ]]; then
        opts="30m 1h 2h 6h 12h 1d 2d 7d"
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
      elif [[ ${prev} == "--program" || ${prev} == "-p" ]]; then
        opts="xyz-hmi xyz-bin-packing xyz-max-hmi-server xyz-studio-max"
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
      else
        opts="--time -t --start-time -s --end-time -e --log-dir -l --upload -u --keep-local-snapshot -k --output-dir -o --program -p --today --yesterday --this-week --skip-version-check --config-dir --simple --interactive -I"
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
      fi
      ;;
    update)
      opts="--force -f --check-only -c --config-dir"
      COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
      ;;
    version)
      opts="--simple -s"
      COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
      ;;
    completion)
      opts="bash zsh fish powershell install"
      COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
      ;;
    *)
      ;;
  esac

  return 0
}

complete -F _logsnap_completion logsnap
