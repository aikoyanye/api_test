package info

import (
	"../log"
	"../net"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

var CsvFile *os.File
var csvWriter *csv.Writer

// 初始化处理csv
func InitCsv(csvInfo *net.CsvInfo){
	if csvInfo.SavePath != ""{
		var err error
		CsvFile, err = os.OpenFile(csvInfo.SavePath, os.O_CREATE|os.O_RDWR, 0644)
		CsvFile.WriteString("\xEF\xBB\xBF")
		if err != nil{
			log.Log.Println("创建CSV文件失败")
			fmt.Printf("[%v]创建CSV文件失败\n", time.Now().Format("2006-01-02 15:04:05"))
			csvInfo.SavePath = ""
			return
		}
		csvWriter = csv.NewWriter(CsvFile)
		csvWriter.Write(csvInfo.Fields)
		csvWriter.Flush()
		log.Log.Println("创建CSV文件成功")
		fmt.Printf("[%v]创建CSV文件成功\n", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// 关闭csv流
func CloseCsv(csvInfo *net.CsvInfo){
	if csvInfo.SavePath != ""{
		CsvFile.Close()
	}
}

// csv数据
func WriteCsv(csvInfo *net.CsvInfo, fields []map[string]string){
	if csvInfo.SavePath != ""{
		InitCsv(csvInfo)
		log.Log.Println("正在将请求结果写入CSV文件...")
		fmt.Printf("[%v]正在将请求结果写入CSV文件...\n", time.Now().Format("2006-01-02 15:04:05"))
		for i, fields := range fields{
			fs := []string{}
			for _, field := range csvInfo.Fields{
				fs = append(fs, fields[field])
			}
			if err := csvWriter.Write(fs); err != nil{
				log.Log.Printf("写入第 %v 条CSV失败：%v\n", i+1, err)
				fmt.Printf("[%v]写入第 %v 条CSV失败：%v\n", i+1, time.Now().Format("2006-01-02 15:04:05"), err)
				continue
			}
			csvWriter.Flush()
			log.Log.Printf("正在写入第 %v 条数据到CSV中...\n", i+1)
			fmt.Printf("[%v]正在写入第 %v 条数据到CSV中...\n", time.Now().Format("2006-01-02 15:04:05"), i+1)
		}
		log.Log.Println("CSV数据写入完毕")
		fmt.Printf("[%v]CSV数据写入完毕\n", time.Now().Format("2006-01-02 15:04:05"))
	}
}