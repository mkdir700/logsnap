package collector

// 处理结果结构体，用于存储每个处理器的处理结果
type ProcessorResult struct {
	processorName string
	outputPath    string
	results       []FileProcessResult
	err           error
}

// FileProcessResult 表示处理结果的统一结构
type FileProcessResult struct {
	FilePath   string // 文件路径
	Err        error  // 错误信息
	FileCount  int    // 处理的文件数量
	FileSize   int64  // 文件大小
	TotalLines int    // 处理的总行数
	MatchLines int    // 匹配的行数
	MatchFiles int    // 匹配的文件数
}


func (p *ProcessorResult) GetTotalLines() int {
	totalLines := 0
	for _, result := range p.results {
		totalLines += result.TotalLines
	}
	return totalLines
}

func (p *ProcessorResult) GetMatchLines() int {
	matchLines := 0
	for _, result := range p.results {
		matchLines += result.MatchLines
	}
	return matchLines
}

func (p *ProcessorResult) GetFileCount() int {
	fileCount := 0
	for _, result := range p.results {
		fileCount += result.FileCount
	}
	return fileCount
}

func (p *ProcessorResult) GetFileSize() int64 {
	fileSize := int64(0)
	for _, result := range p.results {
		fileSize += result.FileSize
	}
	return fileSize
}