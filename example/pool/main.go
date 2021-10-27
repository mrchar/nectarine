package main

import (
	"context"
	"fmt"
	"time"

	util "github.com/mrchar/nectarine"
)

type Runner struct {
	name string
}

func (s Runner) Run(ctx context.Context) error {
	fmt.Printf("现在在执行的是任务:%s\n", s.name)
	time.Sleep(time.Second)
	return nil
}

func main() {
	pool := util.NewPool(10)

	for _, name := range []string{"Task1", "Task2", "Task3"} {
		pool.Add(name, Runner{name: name})
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	pool.Run(ctx)
}
