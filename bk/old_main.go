package old_main

import (
	"./db"
	"./log"
	"./net"
	"./tool"
	"encoding/json"
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var jsonFilePath string
var results net.ServersLice
var logPath string
var resultsMap []map[string]string
var mode string
var h string
var thread_count int

func main(){
	if jsonFilePath != ""{
		log.InitLog(logPath)
		if err := json.Unmarshal(tool.JsonFromFile(jsonFilePath), &results); err != nil{
			log.Log.Printf("JSON文件解析失败，请检查json格式是否规范，程序退出：%v", err)
			panic(fmt.Sprintf("[%v]JSON文件解析失败，请检查json格式是否规范，程序退出：%v",
				time.Now().Format("2006-01-02 15:04:05"), err))
		}
		tool.ReplaceResult(&results)
		log.Log.Printf("共解析出 %v 条请求记录\n", len(results.Servers))
		fmt.Printf("[%v]共解析出 %v 条请求记录\n", time.Now().Format("2006-01-02 15:04:05"), len(results.Servers))
		resultsMap = []map[string]string{}
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
	flag.IntVar(&thread_count, "thcount", runtime.NumCPU(), "默认为当前CPU核心数*2")
	flag.StringVar(&h, "h", "", "有问题去看源码：https://github.com/aikoyanye/api_test")
	flag.Parse()
}

// default默认模式运行
func _default(){
	for _, server := range results.Servers{
		if count, err := strconv.Atoi(server.Count); err == nil{
			for i := 0; i < count; i++ {
				ready(server)
			}
		}else{
			log.Log.Println("Count字段的内容错误，只接受可以转换为整型的字符串")
			fmt.Printf("[%v]Count字段的内容错误，只接受可以转换为整型的字符串\n",
				time.Now().Format("2006-01-02 15:04:05"))
		}
	}
	writeData()
}

// parallel并行模式运行
func _parallel(){
	ch := make(chan int, len(results.Servers))
	for index, server := range results.Servers{
		count, err := strconv.Atoi(server.Count)
		if err != nil{
			log.Log.Println("Count字段的内容错误，只接受可以转换为整型的字符串")
			fmt.Printf("[%v]Count字段的内容错误，只接受可以转换为整型的字符串\n",
				time.Now().Format("2006-01-02 15:04:05"))
			continue
		}
		ch <- index
		go func(){
			for i := 0; i < count; i++ {
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
	var ch chan int
	for _, server := range results.Servers{
		if count, err := strconv.Atoi(server.Count); err == nil{
			ch = make(chan int, count)
			for i := 0; i < count; i++ {
				ch <- i
				go func() {
					ready(server)
					<- ch
				}()
			}
		}else{
			log.Log.Println("Count字段的内容错误，只接受可以转换为整型的字符串")
			fmt.Printf("[%v]Count字段的内容错误，只接受可以转换为整型的字符串\n",
				time.Now().Format("2006-01-02 15:04:05"))
		}
		for true{
			if len(ch) == 0{ break }
			time.Sleep(time.Second)
		}
	}
	chechCh(ch)
}

// all 性能+并行模式
func _all(){
	log.Log.Println("性能并行模式可能会出现不稳定的情况，与电脑配置有关，酌情控制访问数：Count")
	fmt.Printf("[%v]性能并行模式可能会出现不稳定的情况，与电脑配置有关，酌情控制访问数：Count\n",
		time.Now().Format("2006-01-02 15:04:05"))
	count := 0
	for _, server := range results.Servers{
		if c, err := strconv.Atoi(server.Count); err == nil{
			count += c
		}
	}
	ch := make(chan int, count)
	for _, server := range results.Servers{
		if count, err := strconv.Atoi(server.Count); err == nil{
			go func() {
				for i := 0; i < count; i++ {
					ch <- i
					go func() {
						ready(server)
						<- ch
					}()
				}
			}()
		}else{
			log.Log.Println("Count字段的内容错误，只接受可以转换为整型的字符串")
			fmt.Printf("[%v]Count字段的内容错误，只接受可以转换为整型的字符串\n",
				time.Now().Format("2006-01-02 15:04:05"))
		}
	}
	time.Sleep(3 * time.Second)
	chechCh(ch)
}

func writeData(){
	log.Log.Println("此次访问已经结束")
	fmt.Printf("[%v]此次访问已经结束\n", time.Now().Format("2006-01-02 15:04:05"))
	info.WriteCsv(&results.Csv, resultsMap)
	info.WriteSql(&results.Sql, resultsMap)
	info.UploadResult2Api(&results.Upload, resultsMap)
	checkTime()
	defer info.CloseCsv(&results.Csv)
	defer info.CloseDB(&results.Sql)
}

// 返回所有API访问的平均耗时与总耗时
func checkTime(){
	log.Log.Println("正在计算访问总耗时和平均耗时，请稍等...")
	fmt.Printf("[%v]正在计算访问总耗时和平均耗时，请稍等...\n", time.Now().Format("2006-01-02 15:04:05"))
	sum := make(map[string]float64)
	avg := make(map[string]float64)
	count := make(map[string]int)
	for _, m := range resultsMap {
		count[m["RequestApi"]] += 1
		if f, err := strconv.ParseFloat(m["RequestTime"], 64); err == nil{
			sum[m["RequestApi"]] += f
		}else{
			log.Log.Printf("计算访问耗时失败，请自行检查数据：%v\n", err)
			fmt.Printf("[%v]计算访问耗时失败，请自行检查数据：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
			break
		}
	}
	for key, value := range sum{
		avg[key] = value / float64(count[key])
		log.Log.Printf("访问\"%v\"所消耗的时间为：%vs，平均耗时为：%vs", key, value, avg[key])
		fmt.Printf("[%v]访问\"%v\"所消耗的时间为：%vs，平均耗时为：%vs\n",
			time.Now().Format("2006-01-02 15:04:05"), key, value, avg[key])
	}
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
	if success{
		fields := tool.PutResult(resp, result, &server, cost)
		fields["RequestBodyExpected"] = "None"
		if len(server.Expected) != 0{
			fields["RequestBodyExpected"] = tool.EasyCompare(result, server.Expected)
		}
		resultsMap = append(resultsMap, fields)
	}
}
