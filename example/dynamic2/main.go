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
	time.Sleep(3 * time.Second)
	r := executor.GPResult{}
	r.Value = m.param * m.param
	r.Err = nil
	return &r
}

func main() {
	//pool := executor.NewDynamicGPool(context.Background(), 5, 30)
	/*
	   options:
	   executor.WithTaskQueueCap() : set capacity of task queue
	*/
	pool := executor.NewDynamicGPool(context.Background(), 5, 30,
		executor.WithDynamicTaskQueueCap(5),
		executor.WithDetectInterval(time.Second*10),
		executor.WithMeetCondNum(3),
	)
	futureList := make([]executor.Future, 0)
	var f executor.Future
	for i := 0; i < 100; i++ {
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
