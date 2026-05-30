package logger

import "fmt"

type Logger struct {
	verbose bool
}

func New(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
)

func (l *Logger) Simulate(msg string) {
	fmt.Printf("%s[模拟运行]%s %s\n", Yellow, Reset, msg)
}

func (l *Logger) Scan(src string) {
	fmt.Printf("%s[扫描]%s 源目录: %s%s%s\n", Cyan, Reset, Green, src, Reset)
}

func (l *Logger) Depth(limit string) {
	fmt.Printf("%s[深度]%s 最大深度: %s\n", Cyan, Reset, limit)
}

func (l *Logger) Ignore(patterns []string) {
	if len(patterns) > 0 {
		fmt.Printf("%s[忽略]%s %s\n", Cyan, Reset, patterns)
	}
}

func (l *Logger) CreateDir(path string) {
	fmt.Printf("%s[创建]%s %s📁%s %s\n", Green, Reset, Blue, Reset, path)
}

func (l *Logger) CreateFile(path string) {
	fmt.Printf("%s[创建]%s %s📄%s %s %s(空文件，-empty模式)%s\n", Green, Reset, Purple, Reset, path, Yellow, Reset)
}

func (l *Logger) Skip(path string, reason string) {
	fmt.Printf("%s[跳过]%s %s - %s\n", Yellow, Reset, path, reason)
}

func (l *Logger) Complete(dirs, files int) {
	fmt.Printf("%s[完成]%s 共创建 %s%d%s 个目录, %s%d%s 个空文件\n",
		Green, Reset, Cyan, dirs, Reset, Purple, files, Reset)
}

func (l *Logger) Error(msg string) {
	fmt.Printf("%s[错误]%s %s%s%s\n", Red, Reset, White, msg, Reset)
}

func (l *Logger) Warn(msg string) {
	fmt.Printf("%s[警告]%s %s\n", Yellow, Reset, msg)
}

func (l *Logger) Info(msg string) {
	if l.verbose {
		fmt.Printf("%s[信息]%s %s\n", Blue, Reset, msg)
	}
}

func (l *Logger) SimulateMode() {
	fmt.Printf("%s注意: 模拟运行模式，不会实际创建文件/文件夹%s\n", Yellow, Reset)
}

func (l *Logger) PreservedAttrs(path string) {
	if l.verbose {
		fmt.Printf("%s[保留属性]%s %s\n", Blue, Reset, path)
	}
}
