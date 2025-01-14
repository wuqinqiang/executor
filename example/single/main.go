package main

import (
	"context"
	"fmt"
	"github.com/vearne/executor"
	"time"
)

type MyCallable struct {
	param int
}

func (m *MyCallable) Call(ctx context.Context) *executor.GPResult {
	time.Sleep(1 * time.Second)
	r := executor.GPResult{}
	r.Value = m.param * m.param
	r.Err = nil
	return &r
}

func main() {
	pool := executor.NewSingleGPool(context.Background())
	/*
	   options:
	   executor.WithTaskQueueCap() : set capacity of task queue
	*/
	//pool := executor.NewSingleGPool(context.Background(), executor.WithTaskQueueCap(50))
	futureList := make([]executor.Future, 0)
	var f executor.Future
	for i := 0; i < 10; i++ {
		task := &MyCallable{param: i}
		f = pool.Submit(task)
		futureList = append(futureList, f)
	}
	pool.Shutdown() // Prohibit submission of new tasks
	var result *executor.GPResult
	for _, f := range futureList {
		result = f.Get()
		fmt.Println(result.Err, result.Value)
	}
	pool.WaitTerminate()
}
