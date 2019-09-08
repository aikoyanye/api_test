package main

import (
	"fmt"
	"sync/atomic"
)

func main(){
	ch := make(chan int, 6)
	var index int32
	for i := 0; i < 24; i++{
		ch <- i
		go func() {
			for j := 0; j < 1000; j++{
				fmt.Println(atomic.AddInt32(&index, 1))
			}
			<- ch
		}()
	}
	for true{
		if len(ch) == 0{
			fmt.Println("----------------------------------")
			break
		}
	}
}
