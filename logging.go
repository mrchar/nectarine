package nectarine

import "github.com/sirupsen/logrus"

var (
	enableLogger bool
	logger Logger
)

func init() {
	enableLogger = true
	logger = logrus.StandardLogger()
}

type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

// SetLogger 使用指定的logger打印日志
func SetLogger(l Logger){
	logger = l
}

// SetEnableLogger 设置是否打印
func SetEnableLogger(enable bool){
	enableLogger = enable
}