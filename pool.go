package nectarine

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

// Runner 需要执行的动作
type Runner interface {
	// Run Run 被调用执行任务操作，Run不应该阻塞
	Run(ctx context.Context) error
}

// Pool 调度的任务
type Task interface {
	Name() string   // 任务名称
	Replicas() int  // 预期副本数
	Running() int   // 当前副本数
	Runner() Runner // 获取运行者
}

// Pool 运行池
type Pool interface {
	Runner
	// Add 向协程池中添加任务
	Add(name string, runner Runner) error
	// Remove 从协程池中移除任务
	Remove(name string) error
}

func NewPool(size int) Pool {
	return NewSimplePool(size)
}

type SimplePool struct {
	lock      *sync.RWMutex
	semaphore chan struct{}
	runners   map[string]Runner
}

func NewSimplePool(size int) *SimplePool {
	return &SimplePool{
		lock:      &sync.RWMutex{},
		semaphore: make(chan struct{}, size),
		runners:   make(map[string]Runner),
	}
}

func (s *SimplePool) Add(name string, runner Runner) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.runners[name]; ok {
		return errors.New("名称重复")
	}

	s.runners[name] = runner

	return nil
}

func (s *SimplePool) Remove(name string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.runners[name]; !ok {
		return errors.New("任务不存在")
	}

	delete(s.runners, name)

	return nil
}

func (s *SimplePool) Run(ctx context.Context) error {
	for {
		err := s.run(ctx)
		if err != nil {
			break
		}
	}

	return nil
}

func (s *SimplePool) run(ctx context.Context) error {
	s.lock.RLock()
	runners := s.runners
	s.lock.RUnlock()

	for name := range runners {
		runner := runners[name]
		select {
		case <-ctx.Done():
			return errors.New("接收到结束信号")
		default:
			go func() {
				s.semaphore <- struct{}{}
				defer func() { <-s.semaphore }()

				err := runner.Run(ctx)
				if err != nil && enableLogger {
					logger.Println(errors.Wrapf(err, `执行任务"%s失败"`, name))
				}
			}()
		}
	}

	return nil
}
