package nectarine

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// WaitSignal 监听操作系统信号，并阻塞程序执行
func WaitSignal(signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// 如果没有传入信号列表，则使用默认的信号
	if  len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signal.Notify(sigs, signals...)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	if enableLogger {
		logger.Println("正在等待操作系统停止信号...")	
	}
	
	<-done

	if enableLogger {
		logger.Println("接收到操作系统停止信号！")	
	}	
}