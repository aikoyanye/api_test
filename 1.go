package main

import (
	"fmt"
	"strings"
)

func main(){
	s := "[1]"
	s = strings.ReplaceAll(s, "[1]", "23333")
	fmt.Println(s)
}
