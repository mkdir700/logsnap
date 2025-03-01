package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"
)

//go:embed completion/bash-completion.sh
//go:embed completion/zsh-completion.sh
//go:embed completion/fish-completion.fish
//go:embed completion/powershell-completion.ps1
//go:embed completion/install.sh
var completionScripts embed.FS

// completionAction 处理 completion 命令
func completionAction(c *cli.Context) error {
	fmt.Println("请选择 shell 类型：")
	fmt.Println("  logsnap completion bash     - 生成 Bash 自动补全脚本")
	fmt.Println("  logsnap completion zsh      - 生成 Zsh 自动补全脚本")
	fmt.Println("  logsnap completion fish     - 生成 Fish 自动补全脚本")
	fmt.Println("  logsnap completion powershell - 生成 PowerShell 自动补全脚本")
	fmt.Println("  logsnap completion install  - 自动检测 shell 并安装补全脚本")
	return nil
}

// 检测当前使用的 shell
func detectShell() string {
	// 首先尝试从环境变量获取
	shell := os.Getenv("SHELL")

	// 在 Windows 上，尝试检测 PowerShell
	if runtime.GOOS == "windows" || shell == "" {
		// 检查是否在 PowerShell 中运行
		if os.Getenv("PSModulePath") != "" {
			return "powershell"
		}
		// 默认返回 bash，因为它是最常见的
		return "bash"
	}

	// 提取 shell 名称
	shell = filepath.Base(shell)
	return shell
}

// installCompletionAction 自动检测 shell 并安装补全脚本
func installCompletionAction(c *cli.Context) error {
	// 创建临时目录存放脚本文件
	tempDir, err := os.MkdirTemp("", "logsnap-completion")
	if err != nil {
		return fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 提取所有脚本文件到临时目录
	scriptFiles := []string{
		"bash-completion.sh",
		"zsh-completion.sh",
		"fish-completion.fish",
		"powershell-completion.ps1",
		"install.sh",
	}

	for _, fileName := range scriptFiles {
		content, err := completionScripts.ReadFile("completion/" + fileName)
		if err != nil {
			return fmt.Errorf("读取脚本文件 %s 失败: %v", fileName, err)
		}

		err = os.WriteFile(filepath.Join(tempDir, fileName), content, 0644)
		if err != nil {
			return fmt.Errorf("写入临时脚本文件 %s 失败: %v", fileName, err)
		}
	}

	// 设置安装脚本的执行权限
	err = os.Chmod(filepath.Join(tempDir, "install.sh"), 0755)
	if err != nil {
		return fmt.Errorf("设置安装脚本执行权限失败: %v", err)
	}

	// 执行安装脚本
	cmd := exec.Command(filepath.Join(tempDir, "install.sh"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// bashCompletionAction 生成 Bash 自动补全脚本
func bashCompletionAction(c *cli.Context) error {
	content, err := completionScripts.ReadFile("completion/bash-completion.sh")
	if err != nil {
		return fmt.Errorf("读取 Bash 补全脚本失败: %v", err)
	}

	fmt.Println(string(content))
	fmt.Println("\n# 使用方法：")
	fmt.Println("# 方法1：将上述内容保存到文件并加载")
	fmt.Println("#   echo '上述内容' > ~/.logsnap-completion.bash")
	fmt.Println("#   source ~/.logsnap-completion.bash")
	fmt.Println("#")
	fmt.Println("# 方法2：直接加载到当前 shell")
	fmt.Println("#   source <(logsnap completion bash)")
	fmt.Println("#")
	fmt.Println("# 方法3：永久添加到 .bashrc")
	fmt.Println("#   logsnap completion bash >> ~/.bashrc")

	return nil
}

// zshCompletionAction 生成 Zsh 自动补全脚本
func zshCompletionAction(c *cli.Context) error {
	content, err := completionScripts.ReadFile("completion/zsh-completion.sh")
	if err != nil {
		return fmt.Errorf("读取 Zsh 补全脚本失败: %v", err)
	}

	fmt.Println(string(content))
	fmt.Println("\n# 提示：")
	fmt.Println("# 请使用以下命令立即启用补全功能：")
	fmt.Println("#   source <(logsnap completion zsh)")
	fmt.Println("#")
	fmt.Println("# 要永久启用，请将以下行添加到 ~/.zshrc：")
	fmt.Println("#   if type logsnap &>/dev/null; then")
	fmt.Println("#     source <(logsnap completion zsh)")
	fmt.Println("#   fi")

	return nil
}

// fishCompletionAction 生成 Fish 自动补全脚本
func fishCompletionAction(c *cli.Context) error {
	content, err := completionScripts.ReadFile("completion/fish-completion.fish")
	if err != nil {
		return fmt.Errorf("读取 Fish 补全脚本失败: %v", err)
	}

	fmt.Println(string(content))
	fmt.Println("\n# 使用方法：")
	fmt.Println("# 方法1：将上述内容保存到 Fish 补全目录")
	fmt.Println("#   mkdir -p ~/.config/fish/completions")
	fmt.Println("#   echo '上述内容' > ~/.config/fish/completions/logsnap.fish")
	fmt.Println("#")
	fmt.Println("# 方法2：直接加载到当前 shell")
	fmt.Println("#   logsnap completion fish | source")

	return nil
}

// powershellCompletionAction 生成 PowerShell 自动补全脚本
func powershellCompletionAction(c *cli.Context) error {
	content, err := completionScripts.ReadFile("completion/powershell-completion.ps1")
	if err != nil {
		return fmt.Errorf("读取 PowerShell 补全脚本失败: %v", err)
	}

	fmt.Println(string(content))
	fmt.Println("\n# 使用方法：")
	fmt.Println("# 方法1：将上述内容保存到 PowerShell 配置文件")
	fmt.Println("#   echo '上述内容' > $PROFILE.CurrentUserAllHosts")
	fmt.Println("#")
	fmt.Println("# 方法2：直接加载到当前 shell")
	fmt.Println("#   logsnap completion powershell | Out-String | Invoke-Expression")

	return nil
}
