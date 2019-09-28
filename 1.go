package main

import (
	"fmt"
	"strconv"
)

func main(){
	var f float64
	f = 0.0000000002
	s2f, _ := strconv.ParseFloat(strconv.FormatFloat(f, 'f', -1, 64), 64)
	fmt.Println(s2f)
}
