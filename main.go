package main

import (
	"./db"
	"./log"
	"./net"
	"./tool"
	"encoding/json"
	"flag"
	"fmt"
	"time"
)

var jsonFilePath string
var results net.ServersLice
var logPath string
var results_map []map[string]string
var mode string

func main(){
	if jsonFilePath != ""{
		log.InitLog(logPath)
		if err := json.Unmarshal(tool.JsonFromFile(jsonFilePath), &results); err != nil{
			log.Log.Printf("JSON文件解析失败，请检查json格式是否规范，程序退出：%v", err)
			panic(fmt.Sprintf("[%v]JSON文件解析失败，请检查json格式是否规范，程序退出：%v",
				time.Now().Format("2006-01-02 15:04:05"), err))
		}
		log.Log.Printf("共解析出 %v 条请求记录\n", len(results.Servers))
		fmt.Printf("[%v]共解析出 %v 条请求记录\n", time.Now().Format("2006-01-02 15:04:05"), len(results.Servers))
		results_map = []map[string]string{}
		for _, server := range results.Servers{
			ready(server)
		}
	}else{
		flag.Usage()
	}
	log.Log.Println("此次访问已经结束")
	fmt.Printf("[%v]此次访问已经结束\n", time.Now().Format("2006-01-02 15:04:05"))
	info.WriteCsv(&results.Csv, results_map)
	info.WriteSql(&results.Sql, results_map)
	defer info.CloseCsv(&results.Csv)
	defer info.CloseDB(&results.Sql)
}

// 设定程序允许接受的参数
func init(){
	flag.StringVar(&jsonFilePath, "path", "", "json文件的路径，可以是绝对路径，也可以是相对路径")
	flag.StringVar(&logPath, "log", "", "log文件的保存的路径，指定到路径就好，路径后面必须带有斜杠\"/\"")
	flag.StringVar(&mode, "mode", "default", "访问API的模式，有三种模式：default、")
	flag.Parse()
}

// 访问api->比较结果->导出结果->导出log
func ready(server net.Server){
	for i := 0; i < server.Count; i++ {
		result, resp, cost, success := net.HttpGo(&server)
		fields := tool.PutResult(resp, result, &server, cost)
		if success{
			fields["RequestBodyExpected"] = "None"
			if len(server.Expected) != 0{
				fields["RequestBodyExpected"] = tool.EasyCompare(result, server.Expected)
			}
			results_map = append(results_map, fields)
		}
	}
}
