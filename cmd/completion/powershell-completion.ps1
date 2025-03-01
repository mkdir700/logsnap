# logsnap PowerShell completion

Register-ArgumentCompleter -Native -CommandName logsnap -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    $commands = @{
        'collect' = '收集指定时间范围内的日志并上传'
        'update' = '检查并更新程序到最新版本'
        'supported-programs' = '显示支持的程序列表'
        'version' = '显示当前版本信息'
        'completion' = '生成自动补全脚本'
        'help' = '显示帮助信息'
    }
    
    $collectOpts = @(
        '--time', '-t'
        '--start-time', '-s'
        '--end-time', '-e'
        '--log-dir', '-l'
        '--upload', '-u'
        '--keep-local-snapshot', '-k'
        '--output-dir', '-o'
        '--program', '-p'
        '--today'
        '--yesterday'
        '--this-week'
        '--skip-version-check'
        '--config-dir'
        '--simple'
        '--interactive', '-I'
    )
    
    $updateOpts = @(
        '--force', '-f'
        '--check-only', '-c'
        '--config-dir'
    )
    
    $versionOpts = @(
        '--simple', '-s'
    )
    
    $completionOpts = @(
        'bash'
        'zsh'
        'fish'
        'powershell'
        'install'
    )
    
    $commandElements = $commandAst.CommandElements
    $command = $commandElements[0].Value
    
    if ($commandElements.Count -eq 1) {
        return $commands.Keys | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $commands[$_])
        }
    }
    
    $subCommand = $commandElements[1].Value
    
    if ($commandElements.Count -eq 2 -and $subCommand -like "$wordToComplete*") {
        return $commands.Keys | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
            [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $commands[$_])
        }
    }
    
    switch ($subCommand) {
        'collect' {
            return $collectOpts | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        'update' {
            return $updateOpts | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        'version' {
            return $versionOpts | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
        'completion' {
            return $completionOpts | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
                [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
            }
        }
    }
}
