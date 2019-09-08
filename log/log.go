package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Log *log.Logger

func InitLog(path string){
	file, err := os.Create(path + time.Now().Format("2006-01-02 15-04-05") + ".log")
	if err != nil{
		panic(fmt.Sprintf("log文件创建失败，程序退出：%v", err))
	}
	Log = log.New(file, "", log.Ldate|log.Ltime)
}