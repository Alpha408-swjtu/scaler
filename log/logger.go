package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	Logger   = logrus.New()
	LogEntry *logrus.Entry
)

// 日志的输出格式
func Formatter() *nested.Formatter {
	return &nested.Formatter{
		HideKeys:        true,             // 不显示键值对的key
		TimestampFormat: time.RFC3339Nano, // 时间格式
		CallerFirst:     true,             // 调用者信息放在第一位
		NoColors:        true,             // 不显示颜色
		ShowFullLevel:   true,             // 显示完整的日志级别
		CustomCallerFormatter: func(f *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(f.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC()"
			}
			fullPath, line := funcInfo.FileLine(f.PC)
			return fmt.Sprintf(" [%s:%d] ", filepath.Base(fullPath), line)
		}, // 自定义调用者信息，显示代码位置和行号
	}
}

// 日志的初始化
func init() {
	Logger.SetLevel(logrus.DebugLevel)
	// 设置日志格式
	Logger.SetFormatter(Formatter())

	// 设置日志输出到控制台
	Logger.Out = os.Stdout

	// 设置日志钩子，用于显示调用者信息
	Logger.SetReportCaller(true)
	LogEntry = logrus.NewEntry(Logger)
}
