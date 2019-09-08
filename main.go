package main

import (
	"./db"
	"./log"
	"./net"
	"./tool"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
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
		if strings.ToLower(mode) == "default"{
			_default()
		}else if strings.ToLower(mode) == "parallel" {
			_parallel()
		}else if strings.ToLower(mode) == "performance"{
			_performance()
		}else if strings.ToLower(mode) == "all"{
			_all()
		}else{
			log.Log.Println("无效的运行模式，请确认输入的字符是否正确")
			fmt.Printf("[%v]无效的运行模式，请确认输入的字符是否正确\n", time.Now().Format("2006-01-02 15:04:05"))
			flag.Usage()
		}
	}else{
		flag.Usage()
	}
}

// 设定程序允许接受的参数
func init(){
	flag.StringVar(&jsonFilePath, "path", "", "json文件的路径，可以是绝对路径，也可以是相对路径")
	flag.StringVar(&logPath, "log", "", "log文件的保存的路径，指定到路径就好，路径后面必须带有斜杠\"/\"")
	flag.StringVar(&mode, "mode", "default", "访问API的模式，有三种模式：default、parallel、performance、all")
	flag.StringVar(nil, "h", "", "有问题去看源码：https://github.com/aikoyanye/api_test")
	flag.Parse()
}

// default默认模式运行
func _default(){
	for _, server := range results.Servers{
		for i := 0; i < server.Count; i++ {
			ready(server)
		}
	}
	writeData()
}

// parallel并行模式运行
func _parallel(){
	ch := make(chan int, len(results.Servers))
	for index, server := range results.Servers{
		ch <- index
		go func(){
			for i := 0; i < server.Count; i++ {
				ready(server)
			}
			<- ch
		}()
	}
	chechCh(ch)
}

// performance高并发模式
func _performance(){
	log.Log.Println("性能模式可能会出现不稳定的情况，与电脑配置有关，酌情控制访问数：Count")
	fmt.Printf("[%v]性能模式可能会出现不稳定的情况，与电脑配置有关，酌情控制访问数：Count\n",
		time.Now().Format("2006-01-02 15:04:05"))
	ch := make(chan int)
	for _, server := range results.Servers{
		for i := 0; i < server.Count; i++ {
			ch <- i
			go func() {
				ready(server)
				<- ch
			}()
		}
		time.Sleep(time.Second)
	}
	chechCh(ch)
}

// all 性能+并行模式
func _all(){
	log.Log.Println("性能并行模式可能会出现不稳定的情况，与电脑配置有关，酌情控制访问数：Count")
	fmt.Printf("[%v]性能并行模式可能会出现不稳定的情况，与电脑配置有关，酌情控制访问数：Count\n",
		time.Now().Format("2006-01-02 15:04:05"))
	ch := make(chan int)
	for _, server := range results.Servers{
		go func() {
			for i := 0; i < server.Count; i++ {
				ch <- i
				go func() {
					ready(server)
					<- ch
				}()
			}
		}()
	}
	chechCh(ch)
}

func writeData(){
	log.Log.Println("此次访问已经结束")
	fmt.Printf("[%v]此次访问已经结束\n", time.Now().Format("2006-01-02 15:04:05"))
	info.WriteCsv(&results.Csv, results_map)
	info.WriteSql(&results.Sql, results_map)
	defer info.CloseCsv(&results.Csv)
	defer info.CloseDB(&results.Sql)
}

func chechCh(ch chan int){
	for true{
		if len(ch) == 0{
			writeData()
			break
		}
		time.Sleep(3 * time.Second)
	}
}

// 访问api->比较结果->导出结果->导出log
func ready(server net.Server){
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
