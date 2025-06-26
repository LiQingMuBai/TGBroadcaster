package global

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNewGoPool(t *testing.T) {
	// 1. 创建协程池(100个worker)
	pool := NewGoPool(10)
	defer pool.Shutdown()

	// 2. 使用WaitGroup等待所有任务完成
	var wg sync.WaitGroup

	// 3. 提交1000个任务
	for i := 0; i < 100; i++ {
		wg.Add(1)
		taskID := i

		// 使用协程池执行任务
		pool.Go(func() {
			defer wg.Done()
			startTime := time.Now()
			fmt.Printf("Task %d start\n", taskID)
			// 模拟耗时任务
			s := rand.Intn(20)

			//log.Println("s", s)
			time.Sleep(time.Second * time.Duration(s))
			fmt.Printf("Task %d completed\n", taskID)

			fmt.Printf("Task %d  耗时: %v\n", taskID, time.Since(startTime))
		})
	}

	//// 4. 定期打印协程池状态
	//go func() {
	//	ticker := time.NewTicker(500 * time.Millisecond)
	//	defer ticker.Stop()
	//
	//	for range ticker.C {
	//		active, total := pool.Stats()
	//		fmt.Printf("Pool stats - Active: %d, Total: %d, Pending: %d\n",
	//			active, total, pool.PendingTasks())
	//	}
	//}()

	// 5. 等待所有任务完成
	wg.Wait()
	fmt.Println("All tasks completed")
}
