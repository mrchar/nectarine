package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	util "github.com/mrchar/nectarine"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		fmt.Println("启动任务！")

		for{
			select {
				case <-ctx.Done():
					return
				default:
					fmt.Printf("任务执行中，当前时间： %s\n",time.Now().Format(time.RFC3339))
					// 使用Sleep模拟正在执行的任务
					time.Sleep(5* time.Second)
			}
		}
	}(ctx)

	util.WaitSignal()

	fmt.Println("接收到系统信号，正在停止程序...")
	cancel()

	wg.Wait()
	fmt.Println("所有任务结束，程序退出！")
}