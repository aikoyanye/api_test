package tool

import (
	"../net"
)

type ThreadPool struct {
	Job 	chan net.Server
	Result 	chan int
	Run		func(server net.Server)
	Count 	int
}

func GetThreadPool(count int, run func(server net.Server)) *ThreadPool {
	return &ThreadPool{Count: count, Job: make(chan net.Server, count), Run: run, Result: make(chan int, count)}
}

func (tp ThreadPool) Start(){
	for i := 0; i < tp.Count; i++{
		go func() {
			for {
				tp.Result <- 1
				server := <-tp.Job
				<- tp.Result
				tp.Run(server)
			}
		}()
	}
}

func (tp ThreadPool) Stop(){
	close(tp.Job)
}

func (tp ThreadPool) AddTask(server net.Server){
	tp.Job <- server
}

func (tp ThreadPool) SetRun(f func(server net.Server)){
	tp.Run = f
}

func (tp ThreadPool) GetCurrentJob() int {
	return len(tp.Job)
}
