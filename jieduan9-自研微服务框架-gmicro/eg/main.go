package main

import (
	"context"
	"fmt"
	"time"

	"errors"

	"golang.org/x/sync/errgroup"
)

func main() {
	//errgroup的go方法内部会启动一个goroutine
	eg, ctx := errgroup.WithContext(context.Background())

	// 用法
	eg.Go(func() error {
		fmt.Println("doing task1")
		time.Sleep(5 * time.Second)
		return errors.New("task1 error")
	})

	eg.Go(func() error {
		for {
			select {
			// 	time.After (时间) = 定时触发的通道，它返回一个 只读 channel
			// 等待你指定的时间（这里是 1 秒）--- time.Second就是1秒的意思
			// 时间到了 → 自动往 channel 里发一个值
			// 没到时间 → 一直阻塞
			case <-time.After(time.Second):
				fmt.Println("doing task2")
			case <-ctx.Done():
				fmt.Println("task2 canceled")
				return ctx.Err()
			}
		}
	})

	eg.Go(func() error {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("doing task3")
			case <-ctx.Done():
				fmt.Println("task3 canceled")
				return ctx.Err()
			}
		}
	})
	// eg.Wait () 是等【服务结束】 与 sync.WaitGroup (wg.Add/Done/Wait) 一般可以搭配使用
	// 	用来等协程结束
	// 只要有一个服务挂了，全部取消
	// 最后主 goroutine 调用 eg.Wait() 一直阻塞到所有服务都停了
	err := eg.Wait()
	if err != nil {
		fmt.Println("task failed")
	} else {
		fmt.Println("task success")
	}
}
