package ui

import (
	"fmt"
	"strings"
	"time"

	"logsnap/remote"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// 版本信息结构
type VersionInfo struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseDate    time.Time
	ReleaseNotes   string
	DownloadURL    string
	ForceUpdate    bool
	UpdateMessage  string
}

// 版本检查结果消息
type versionCheckResultMsg struct {
	Info      VersionInfo
	HasUpdate bool
	Error     error
}

// 下载进度消息
type downloadProgressMsg struct {
	Progress float64 // 0-1 之间的进度
	Done     bool
	Error    error
	FilePath string
}

// 安装结果消息
type installResultMsg struct {
	Success bool
	Error   error
}

// 初始化版本检查命令
func checkVersionCmd(configManager *remote.ConfigManager) tea.Cmd {
	return func() tea.Msg {
		// 使用 ConfigManager 检查更新
		hasUpdate, latestVersion, downloadURL, forceUpdate, updateMessage, err := configManager.CheckForUpdates()
		if err != nil {
			return versionCheckResultMsg{
				Error: err,
			}
		}

		// 获取当前版本
		currentVersion := configManager.GetLocalConfig().GetVersion()

		// 创建版本信息
		info := VersionInfo{
			CurrentVersion: currentVersion,
			LatestVersion:  latestVersion,
			ReleaseDate:    time.Now(), // 远程配置中没有提供发布日期，使用当前时间
			ReleaseNotes:   updateMessage,
			DownloadURL:    downloadURL,
			ForceUpdate:    forceUpdate,
			UpdateMessage:  updateMessage,
		}

		return versionCheckResultMsg{
			Info:      info,
			HasUpdate: hasUpdate,
			Error:     nil,
		}
	}
}

// 下载更新命令
func downloadUpdateCmd(configManager *remote.ConfigManager, downloadURL string) tea.Cmd {
	return func() tea.Msg {
		// 使用 ConfigManager 下载更新
		filePath, err := configManager.DownloadUpdate(downloadURL)
		if err != nil {
			return downloadProgressMsg{
				Done:  true,
				Error: err,
			}
		}

		return downloadProgressMsg{
			Progress: 1.0,
			Done:     true,
			FilePath: filePath,
		}
	}
}

// 安装更新命令
func installUpdateCmd(configManager *remote.ConfigManager, updateFilePath string) tea.Cmd {
	return func() tea.Msg {
		// 使用 ConfigManager 安装更新
		err := configManager.InstallUpdate(updateFilePath)
		if err != nil {
			return installResultMsg{
				Success: false,
				Error:   err,
			}
		}

		return installResultMsg{
			Success: true,
		}
	}
}

// 版本检查组件模型
type VersionCheckerModel struct {
	spinner        spinner.Model
	configManager  *remote.ConfigManager
	currentVersion string
	versionInfo    VersionInfo
	checking       bool
	hasUpdate      bool
	error          error
	expanded       bool // 是否展开显示详细信息

	// 更新相关状态
	downloading      bool
	downloadProgress float64
	downloadError    error
	downloadFilePath string

	installing     bool
	installError   error
	installSuccess bool
}

// 创建新的版本检查组件
func NewVersionChecker(configManager *remote.ConfigManager) VersionCheckerModel {
	s := spinner.New()
	s.Style = SpinnerStyle

	// 获取当前版本
	currentVersion := configManager.GetLocalConfig().GetVersion()

	return VersionCheckerModel{
		spinner:        s,
		configManager:  configManager,
		currentVersion: currentVersion,
		checking:       false,
		hasUpdate:      false,
		expanded:       false,
		downloading:    false,
		installing:     false,
	}
}

// 初始化组件
func (m VersionCheckerModel) Init() tea.Cmd {
	return m.CheckVersion()
}

// 开始检查版本
func (m VersionCheckerModel) CheckVersion() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		checkVersionCmd(m.configManager),
	)
}

// 开始下载更新
func (m VersionCheckerModel) DownloadUpdate() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		downloadUpdateCmd(m.configManager, m.versionInfo.DownloadURL),
	)
}

// 开始安装更新
func (m VersionCheckerModel) InstallUpdate() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		installUpdateCmd(m.configManager, m.downloadFilePath),
	)
}

// 倒计时退出命令
func countdownExitCmd() tea.Cmd {
	// 不再使用倒计时，直接退出
	return tea.Quit
}

// 倒计时消息
type countdownMsg int

// 更新组件状态
func (m VersionCheckerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			// 按 c 键检查更新
			if !m.checking && !m.downloading && !m.installing {
				m.checking = true
				return m, m.CheckVersion()
			}
		case "e":
			// 按 e 键展开/折叠详情
			if m.hasUpdate {
				m.expanded = !m.expanded
			}
		case "d":
			// 按 d 键下载更新
			if m.hasUpdate && !m.downloading && !m.installing {
				m.downloading = true
				m.downloadProgress = 0
				m.downloadError = nil
				return m, m.DownloadUpdate()
			}
		case "i":
			// 按 i 键安装更新
			if m.downloadFilePath != "" && !m.installing {
				m.installing = true
				m.installError = nil
				return m, m.InstallUpdate()
			}
		case "q", "ctrl+c":
			// 按 q 或 ctrl+c 立即退出
			return m, tea.Quit
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case versionCheckResultMsg:
		m.checking = false
		if msg.Error != nil {
			m.error = msg.Error
		} else {
			m.versionInfo = msg.Info
			m.hasUpdate = msg.HasUpdate
			m.error = nil

			// 如果没有更新且启用了自动退出，开始倒计时退出
			if !m.hasUpdate {
				return m, countdownExitCmd()
			}
		}
		return m, nil

	case downloadProgressMsg:
		if msg.Done {
			m.downloading = false
			m.downloadError = msg.Error
			m.downloadFilePath = msg.FilePath
		} else {
			m.downloadProgress = msg.Progress
		}
		return m, nil

	case installResultMsg:
		m.installing = false
		if !msg.Success {
			m.installError = msg.Error
		} else {
			m.installSuccess = true

			return m, countdownExitCmd()
		}
		return m, nil

	case countdownMsg:
		// 处理倒计时消息
		remaining := int(msg)

		if remaining <= 0 {
			// 倒计时结束，退出程序
			return m, tea.Quit
		}

		// 继续倒计时
		return m, countdownExitCmd()
	}

	return m, nil
}

// 渲染组件
func (m VersionCheckerModel) View() string {
	var sb strings.Builder

	// 如果正在检查
	if m.checking {
		sb.WriteString(m.spinner.View() + " 正在检查更新...\n")
		return sb.String()
	}

	// 如果有错误
	if m.error != nil {
		sb.WriteString("检查更新失败: " + m.error.Error() + "\n")
		sb.WriteString("按 'c' 重新检查 | 按 'q' 退出")
		return sb.String()
	}

	// 如果已检查但没有更新
	if m.versionInfo.CurrentVersion != "" && !m.hasUpdate {
		sb.WriteString("✓ 已是最新版本\n")
		sb.WriteString("当前版本: " + m.versionInfo.CurrentVersion + "\n")
		sb.WriteString("按 'c' 重新检查 | 按 'q' 退出")
		return sb.String()
	}

	// 如果有更新
	if m.hasUpdate {
		sb.WriteString("⚠ 发现新版本\n")
		sb.WriteString("当前版本: " + m.versionInfo.CurrentVersion + "\n")
		sb.WriteString("最新版本: " + m.versionInfo.LatestVersion + "\n")

		// 如果有强制更新标志
		if m.versionInfo.ForceUpdate {
			sb.WriteString("⚠ 此更新为强制更新，请尽快更新\n")
		}

		// 如果展开显示详情
		if m.expanded && m.versionInfo.UpdateMessage != "" {
			sb.WriteString("更新内容:\n")
			sb.WriteString(m.versionInfo.UpdateMessage + "\n")
		}

		// 如果正在下载
		if m.downloading {
			sb.WriteString(m.spinner.View() + " 正在下载更新...\n")
			return sb.String()
		}

		// 如果下载出错
		if m.downloadError != nil {
			sb.WriteString("下载更新失败: " + m.downloadError.Error() + "\n")
		}

		// 如果下载完成但未安装
		if m.downloadFilePath != "" && !m.installing && !m.installSuccess {
			sb.WriteString("✓ 下载完成\n")

			// 如果正在安装
			if m.installing {
				sb.WriteString(m.spinner.View() + " 正在安装更新...\n")
				return sb.String()
			}

			// 显示安装选项
			sb.WriteString("按 'i' 安装更新\n")
		}

		// 如果安装出错
		if m.installError != nil {
			sb.WriteString("安装更新失败: " + m.installError.Error() + "\n")
		}

		// 如果安装成功
		if m.installSuccess {
			sb.WriteString("✓ 更新安装成功！请重新启动程序以应用更新。\n")
		}

		// 根据当前状态显示不同的帮助信息
		if m.downloadFilePath == "" && !m.downloading {
			sb.WriteString("按 'e' 查看详情 | 按 'd' 下载更新 | 按 'c' 重新检查 | 按 'q' 退出")
		} else if !m.installSuccess && !m.installing && m.downloadFilePath != "" {
			sb.WriteString("按 'i' 安装更新 | 按 'c' 重新检查 | 按 'q' 退出")
		} else if m.installSuccess {
			sb.WriteString("按 'q' 退出")
		} else {
			sb.WriteString("按 'c' 重新检查 | 按 'q' 退出")
		}

		return sb.String()
	}

	// 初始状态
	sb.WriteString("按 'c' 检查更新 | 按 'q' 退出")
	return sb.String()
}

// 作为独立程序运行版本检查器
func RunVersionChecker(configManager *remote.ConfigManager) {
	p := tea.NewProgram(NewVersionChecker(configManager))
	if _, err := p.Run(); err != nil {
		fmt.Println("运行版本检查器时出错:", err)
	}
}

// AutoUpdateVersionCheckerModel 是一个自动执行更新流程的版本检查器模型
type AutoUpdateVersionCheckerModel struct {
	VersionCheckerModel
}

// 创建新的自动更新版本检查器
func NewAutoUpdateVersionChecker(configManager *remote.ConfigManager) AutoUpdateVersionCheckerModel {
	return AutoUpdateVersionCheckerModel{
		VersionCheckerModel: NewVersionChecker(configManager),
	}
}

// 更新组件状态，重写以支持自动更新流程
func (m AutoUpdateVersionCheckerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case versionCheckResultMsg:
		m.checking = false
		if msg.Error != nil {
			m.error = msg.Error
			// 如果检查出错，立即退出
			return m, tea.Quit
		} else {
			m.versionInfo = msg.Info
			m.hasUpdate = msg.HasUpdate
			m.error = nil

			// 如果没有更新，立即退出
			if !m.hasUpdate {
				return m, tea.Quit
			}

			// 如果有更新，自动开始下载
			m.downloading = true
			m.downloadProgress = 0
			m.downloadError = nil
			return m, m.DownloadUpdate()
		}

	case downloadProgressMsg:
		if msg.Done {
			m.downloading = false
			m.downloadError = msg.Error

			// 如果下载出错，立即退出
			if msg.Error != nil {
				return m, tea.Quit
			}

			m.downloadFilePath = msg.FilePath

			// 下载完成后自动开始安装
			m.installing = true
			m.installError = nil
			return m, m.InstallUpdate()
		} else {
			m.downloadProgress = msg.Progress
		}
		return m, nil

	case installResultMsg:
		m.installing = false
		if !msg.Success {
			m.installError = msg.Error
			// 如果安装失败，立即退出
			return m, tea.Quit
		} else {
			m.installSuccess = true
			// 安装成功后立即退出
			return m, tea.Quit
		}
	}

	// 处理其他消息类型
	model, cmd := m.VersionCheckerModel.Update(msg)
	if updatedModel, ok := model.(VersionCheckerModel); ok {
		m.VersionCheckerModel = updatedModel
		return m, cmd
	}
	return m, cmd
}

// 渲染组件
func (m AutoUpdateVersionCheckerModel) View() string {
	var sb strings.Builder

	// 如果正在检查
	if m.checking {
		sb.WriteString(m.spinner.View() + " 正在检查更新...\n")
		return sb.String()
	}

	// 如果有错误
	if m.error != nil {
		sb.WriteString("检查更新失败: " + m.error.Error() + "\n")
		sb.WriteString("按 'q' 退出")
		return sb.String()
	}

	// 如果已检查但没有更新
	if m.versionInfo.CurrentVersion != "" && !m.hasUpdate {
		sb.WriteString("✓ 已是最新版本\n")
		sb.WriteString("当前版本: " + m.versionInfo.CurrentVersion + "\n")
		sb.WriteString("按 'q' 退出")
		return sb.String()
	}

	// 如果有更新
	if m.hasUpdate {
		sb.WriteString("⚠ 发现新版本\n")
		sb.WriteString("当前版本: " + m.versionInfo.CurrentVersion + "\n")
		sb.WriteString("最新版本: " + m.versionInfo.LatestVersion + "\n")

		// 如果有强制更新标志
		if m.versionInfo.ForceUpdate {
			sb.WriteString("⚠ 此更新为强制更新\n")
		}

		// 如果正在下载
		if m.downloading {
			sb.WriteString(m.spinner.View() + " 正在自动下载更新...\n")
			return sb.String()
		}

		// 如果下载出错
		if m.downloadError != nil {
			sb.WriteString("下载更新失败: " + m.downloadError.Error() + "\n")
			sb.WriteString("按 'q' 退出")
			return sb.String()
		}

		// 如果正在安装
		if m.installing {
			sb.WriteString(m.spinner.View() + " 正在自动安装更新...\n")
			return sb.String()
		}

		// 如果安装出错
		if m.installError != nil {
			sb.WriteString("安装更新失败: " + m.installError.Error() + "\n")
			sb.WriteString("按 'q' 退出")
			return sb.String()
		}

		// 如果安装成功
		if m.installSuccess {
			sb.WriteString("✓ 更新安装成功，请重新启动程序以应用更新！\n")
			sb.WriteString("按 'q' 退出")
		}

		return sb.String()
	}

	// 初始状态
	sb.WriteString("正在准备检查更新...\n")
	return sb.String()
}

// 运行自动更新版本检查器
func RunAutoUpdateVersionChecker(configManager *remote.ConfigManager) {
	p := tea.NewProgram(NewAutoUpdateVersionChecker(configManager))
	if _, err := p.Run(); err != nil {
		fmt.Println("运行自动更新版本检查器时出错:", err)
	}
}
