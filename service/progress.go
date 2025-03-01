package service

// ProgressCallback 进度回调函数类型
type ProgressCallback func(stage string, progress int, message string)

// ProgressReporter 进度报告接口
type ProgressReporter interface {
	Report(stage string, percentage int, message string)
}

// DefaultProgressReporter 默认进度报告实现
type DefaultProgressReporter struct {
	callback ProgressCallback
}

// Report 报告进度
func (r *DefaultProgressReporter) Report(stage string, percentage int, message string) {
	if r.callback != nil {
		r.callback(stage, percentage, message)
	}
}

// NewProgressReporter 创建新的进度报告器
func NewProgressReporter(callback ProgressCallback) ProgressReporter {
	return &DefaultProgressReporter{callback: callback}
}
