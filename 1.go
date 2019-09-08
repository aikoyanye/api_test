package main

import (
	"encoding/csv"
	"os"
)

func main(){
	file, _ := os.OpenFile("2.csv", os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	w.Write([]string{"1erwer", "2werwer", "3werwerwe"})
	w.Flush()
	w.Write([]string{"4", "5", "6"})
	w.Flush()
	w.Write([]string{"7", "8", "9"})
	w.Flush()
}
