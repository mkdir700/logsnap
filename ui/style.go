package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// 颜色定义
var (
	// 主题颜色
	PrimaryColor   = lipgloss.Color("69")
	SecondaryColor = lipgloss.Color("39")
	AccentColor    = lipgloss.Color("168")

	// 文本颜色
	TextColor       = lipgloss.Color("252")
	SubtleTextColor = lipgloss.Color("241")
	ErrorColor      = lipgloss.Color("196")
	SuccessColor    = lipgloss.Color("76")
	WarningColor    = lipgloss.Color("208")

	// 背景颜色
	BackgroundColor = lipgloss.Color("236")
	HighlightColor  = lipgloss.Color("237")
)

// 基础样式
var (
	// 基础文本样式
	BaseStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	// 标题样式
	TitleStyle = BaseStyle.Copy().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	// 副标题样式
	SubtitleStyle = BaseStyle.Copy().
			Foreground(SecondaryColor).
			Bold(true)

	// 强调文本样式
	EmphasisStyle = BaseStyle.Copy().
			Foreground(AccentColor).
			Italic(true)

	// 帮助文本样式
	HelpStyle = BaseStyle.Copy().
			Foreground(SubtleTextColor).
			Margin(1, 0)

	// 错误文本样式
	ErrorStyle = BaseStyle.Copy().
			Foreground(ErrorColor)

	// 成功文本样式
	SuccessStyle = BaseStyle.Copy().
			Foreground(SuccessColor)

	// 警告文本样式
	WarningStyle = BaseStyle.Copy().
			Foreground(WarningColor)
)

// 组件样式
var (
	// 容器样式
	ContainerStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(1, 2)

	// 边框容器样式
	BorderedContainerStyle = ContainerStyle.Copy().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(PrimaryColor)

	// 加载指示器样式
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor)

	// 按钮样式
	ButtonStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(PrimaryColor).
			Padding(0, 3).
			Margin(0, 1).
			Bold(true)

	// 选中项样式
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(SecondaryColor).
				Padding(0, 1)

	// 进度条样式
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor)

	// 表格样式
	TableHeaderStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(PrimaryColor).
				Bold(true).
				Padding(0, 1).
				Align(lipgloss.Center)

	TableCellStyle = lipgloss.NewStyle().
			Padding(0, 1)
)

// 布局样式
var (
	// 应用程序主容器样式
	AppStyle = lipgloss.NewStyle().
			Margin(1, 2, 0, 2)

	// 水平分割线样式
	DividerStyle = lipgloss.NewStyle().
			Foreground(SubtleTextColor).
			Render("─────────────────────────────────")

	// 点样式（用于加载指示等）
	DotStyle = HelpStyle.Copy().
			UnsetMargins()

	// 版本信息样式 - 确保左对齐
	VersionInfoStyle = BaseStyle.Copy().
				Align(lipgloss.Left).
				UnsetMargins().
				UnsetPadding()
)

// 创建自定义宽度的文本样式
func NewWidthStyle(width int) lipgloss.Style {
	return BaseStyle.Copy().Width(width)
}

// 创建自定义高度的文本样式
func NewHeightStyle(height int) lipgloss.Style {
	return BaseStyle.Copy().Height(height)
}

// 创建自定义尺寸的文本样式
func NewSizeStyle(width, height int) lipgloss.Style {
	return BaseStyle.Copy().Width(width).Height(height)
}

// 创建自定义对齐方式的文本样式
func NewAlignStyle(align lipgloss.Position) lipgloss.Style {
	return BaseStyle.Copy().Align(align)
}

// 创建自定义边框样式
func NewBorderStyle(borderStyle lipgloss.Border, borderColor lipgloss.Color) lipgloss.Style {
	return ContainerStyle.Copy().
		Border(borderStyle).
		BorderForeground(borderColor)
}
