package ui

import (
	"fmt"
	"runtime"
	"strings"

	"logsnap/remote"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ShowVersionModel 是显示版本信息的模型
type ShowVersionModel struct {
	spinner       spinner.Model
	configManager *remote.ConfigManager
	versionInfo   VersionInfo
	checking      bool
	hasUpdate     bool
	error         error
	quitting      bool
}

// NewShowVersionModel 创建一个新的版本显示模型
func NewShowVersionModel(configManager *remote.ConfigManager) ShowVersionModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return ShowVersionModel{
		spinner:       s,
		configManager: configManager,
		checking:      true,
	}
}

// Init 初始化模型
func (m ShowVersionModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		checkVersionCmd(m.configManager),
	)
}

// Update 更新模型状态
func (m ShowVersionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc || msg.Type == tea.KeyEnter {
			return m, tea.Quit
		}
		return m, nil

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
		}
		// 检查完成后立即退出
		return m, tea.Quit

	default:
		return m, nil
	}
}

// View 渲染视图
func (m ShowVersionModel) View() string {
	var s strings.Builder

	// 样式定义
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)

	// 信息样式 - 保持一致的标签宽度
	infoLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(10).
		Align(lipgloss.Right)

	infoValueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// 更新提示样式
	updateStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("43"))

	// 错误样式
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196"))

	// 注释样式
	noteStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// 添加标题
	s.WriteString(titleStyle.Render("LogSnap 版本信息") + "\n\n")

	if m.checking {
		s.WriteString(fmt.Sprintf("%s 正在检查版本信息...", m.spinner.View()))
		return s.String()
	}

	// 创建信息行
	infoLine := func(label, value string) string {
		return lipgloss.JoinHorizontal(lipgloss.Top,
			infoLabelStyle.Render(label+": "),
			infoValueStyle.Render(value),
		)
	}

	// 基本信息部分
	if m.error != nil {
		s.WriteString(infoLine("当前版本", m.configManager.GetLocalConfig().GetVersion()) + "\n")
		s.WriteString(infoLine("系统", runtime.GOOS) + "\n")
		s.WriteString(infoLine("架构", runtime.GOARCH) + "\n\n")
		s.WriteString(errorStyle.Render(fmt.Sprintf("检查更新失败: %v", m.error)) + "\n")
	} else {
		s.WriteString(infoLine("当前版本", m.versionInfo.CurrentVersion) + "\n")
		s.WriteString(infoLine("系统", runtime.GOOS) + "\n")
		s.WriteString(infoLine("架构", runtime.GOARCH) + "\n\n")

		if m.hasUpdate {
			s.WriteString(updateStyle.Render(fmt.Sprintf("发现新版本: %s", m.versionInfo.LatestVersion)) + "\n")
			s.WriteString(updateStyle.Render("可以使用 'logsnap update' 命令更新") + "\n")

			if m.versionInfo.UpdateMessage != "" {
				s.WriteString("\n" + noteStyle.Render("更新说明:") + "\n")
				s.WriteString(noteStyle.Render(m.versionInfo.UpdateMessage) + "\n")
			}
		} else {
			s.WriteString(noteStyle.Render("当前已是最新版本") + "\n")
		}
	}

	return s.String()
}

// RunShowVersion 运行版本显示界面
func RunShowVersion(configManager *remote.ConfigManager) {
	p := tea.NewProgram(NewShowVersionModel(configManager))
	if _, err := p.Run(); err != nil {
		fmt.Printf("运行版本显示界面出错: %v\n", err)
	}
}
