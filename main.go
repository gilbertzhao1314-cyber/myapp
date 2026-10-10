package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// ctx, cancel := context.WithCancel(context.Background())
	// for i := 1; i <= 3; i++ {
	// 	go worker(ctx, i)
	// }
	// time.Sleep(2 * time.Second)
	// cancel() // 通知所有 worker 退出
	// time.Sleep(100 * time.Millisecond)ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel() // 一定要 defer cancel，释放资源

	// if err := slowOp(ctx); err != nil {
	// 	fmt.Println(err)                                      // context deadline exceeded
	// 	fmt.Println(errors.Is(err, context.DeadlineExceeded)) // true
	// }
	var wg sync.WaitGroup
	counter := SafeCounter{}

	// 启动 1000 个协程并发执行加法
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go counter.Increment(&wg)
	}

	wg.Wait()
	fmt.Println("Final Count:", counter.count) // 永远保证输出 1000
}

func worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d 退出: %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("worker %d 工作中\n", id)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func slowOp(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second): // 模拟耗时操作
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment(wg *sync.WaitGroup) {
	defer wg.Done()

	// 1. 在修改数据前加锁
	c.mu.Lock()

	// 3. 使用 defer 确保在函数退出时一定会释放锁，防止死锁
	defer c.mu.Unlock()

	// 2. 临界区：此时只有当前协程能修改 count
	c.count++
}
