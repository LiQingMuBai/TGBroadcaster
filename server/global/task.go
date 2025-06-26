package global

import (
	"sync"
	"sync/atomic"
	"time"
)

// TaskFunc 表示要执行的任务函数类型
type TaskFunc func()

// GoPool 固定大小协程池
type GoPool struct {
	taskChan     chan TaskFunc  // 任务通道
	workerCount  int            // 协程数量
	wg           sync.WaitGroup // 等待组
	activeCount  atomic.Int32   // 活跃worker计数
	totalTasks   atomic.Int64   // 总任务计数
	shutdownChan chan struct{}  // 关闭信号
}

// NewGoPool 创建新的协程池
func NewGoPool(workerCount int) *GoPool {
	pool := &GoPool{
		taskChan:     make(chan TaskFunc, workerCount*2), // 缓冲大小为worker数量的两倍
		workerCount:  workerCount,
		shutdownChan: make(chan struct{}),
	}

	// 启动worker
	pool.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go pool.worker()
	}

	return pool
}

// worker 工作协程
func (p *GoPool) worker() {
	defer p.wg.Done()
	p.activeCount.Add(1)

	for {
		select {
		case task, ok := <-p.taskChan:
			if !ok {
				p.activeCount.Add(-1)
				return // 通道关闭，退出协程
			}

			// 执行任务
			task()
			p.totalTasks.Add(1)

		case <-p.shutdownChan:
			p.activeCount.Add(-1)
			return // 收到关闭信号，退出协程
		}
	}
}

// Go 替代go func()的池化版本
func (p *GoPool) Go(f TaskFunc) {
	p.taskChan <- f
}

// GoWithTimeout 带超时的任务提交
func (p *GoPool) GoWithTimeout(f TaskFunc, timeout time.Duration) bool {
	select {
	case p.taskChan <- f:
		return true
	case <-time.After(timeout):
		return false
	case <-p.shutdownChan:
		return false
	}
}

// Shutdown 优雅关闭协程池
func (p *GoPool) Shutdown() {
	close(p.taskChan)
	close(p.shutdownChan)
	p.wg.Wait()
}

// Stats 获取协程池统计信息
func (p *GoPool) Stats() (active int32, total int64) {
	return p.activeCount.Load(), p.totalTasks.Load()
}

// PendingTasks 获取待处理任务数
func (p *GoPool) PendingTasks() int {
	return len(p.taskChan)
}
