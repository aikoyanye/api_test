package tool

import (
	"../log"
	"../net"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var Fields = []string{
	"ResponseTransferEncoding",
	"ResponseContentLength",
	"ResponseStatus",
	"ResponseStatusCode",
	"ResponseProto",
	"ResponseHeader",
	"ResponseBody",
	"RequestApi",
	"RequestMethod",
	"RequestHeader",
	"RequestBody",
	"RequestTime",
	"RequestBodyExpected",
}

// str转map[string]string
func Str2map(str string) map[string]string {
	var result map[string]string
	if err := json.Unmarshal([]byte(strings.ReplaceAll(str, "'", "\"")), &result); err != nil{
		log.Log.Printf("读取json文件时发生错误，请检查json格式是否正确：%v", err)
		panic(fmt.Sprintf("[%v]读取json文件时发生错误，请检查json格式是否正确：%v",
			time.Now().Format("2006-01-02 15:04:05"), err))
	}
	return result
}

// 对比map
func EasyCompareMap(format, result map[string]interface{}, k string) bool {
	re := true
	if len(format) != len(result){
		re = false
		log.Log.Printf("预期结果与返回结果\"%v\"的字段数目不一致\n", k)
		fmt.Printf("[%v]预期结果与返回结果\"%v\"的字段数目不一致\n", time.Now().Format("2006-01-02 15:04:05"), k)
	}
	for key, value := range format{
		if data, err := result[key]; err {
			if reflect.TypeOf(data) == reflect.TypeOf(value){
				log.Log.Printf("\"%v\"看起来是符合预期的\n", key)
				fmt.Printf("[%v]\"%v\"看起来是符合预期的\n", time.Now().Format("2006-01-02 15:04:05"), key)
				// 判断结果中的map属性，但是不能作为判断依据，所以注释了
				//switch data.(type) {
				//case map[string]interface{}:
				//	EasyCompareMap(value.(map[string]interface{}), data.(map[string]interface{}), key)
				//}
			}else{
				re = false
				log.Log.Printf("\"%v\"看起来是不符合预期的\n", key)
				fmt.Printf("[%v]\"%v\"看起来是不符合预期的\n", time.Now().Format("2006-01-02 15:04:05"), key)
			}
		}else{
			re = false
			log.Log.Printf("返回结果中不存在字段: %v\n", key)
			fmt.Printf("[%v]返回结果中不存在字段: %v\n", time.Now().Format("2006-01-02 15:04:05"), key)
		}
	}
	return re
}

// 对比期望结果
func EasyCompare(format, result map[string]interface{}) string {
	re := true
	reStr := "Unknown"
	log.Log.Println("正在比较请求返回的数据")
	fmt.Printf("[%v]正在比较请求返回的数据\n", time.Now().Format("2006-01-02 15:04:05"))
	re = EasyCompareMap(format, result, "root")
	if re{
		log.Log.Println("返回结果貌似与预期结果不一致")
		fmt.Printf("[%v]返回结果貌似与预期结果不一致\n", time.Now().Format("2006-01-02 15:04:05"))
		reStr = "Unconformity"
	}else{
		log.Log.Println("返回结果貌似与预期结果一致")
		fmt.Printf("[%v]返回结果貌似与预期结果一致\n", time.Now().Format("2006-01-02 15:04:05"))
		reStr = "Conformity"
	}
	return reStr
}

// 从文件中读取json数据
func JsonFromFile(path string) []byte {
	if file, err := os.Open(path); err != nil{
		log.Log.Printf("JSON文件打开失败：%v", err)
		panic(fmt.Sprintf("[%v]JSON文件打开失败：%v", time.Now().Format("2006-01-02 15:04:05"), err))
	}else{
		if result, err := ioutil.ReadAll(file); err == nil{
			return result
		}
	}
	return nil
}

func PutResult(resp *http.Response, m map[string]interface{}, server *net.Server, cost float64) map[string]string {
	fields := make(map[string]string)
	fields["ResponseTransferEncoding"] = strings.Join(resp.TransferEncoding, ";")
	fields["ResponseContentLength"] = string(resp.ContentLength)
	fields["ResponseStatus"] = resp.Status
	fields["ResponseStatusCode"] = string(resp.StatusCode)
	fields["ResponseProto"] = resp.Proto
	header := ""
	for key, value := range resp.Header{
		header += fmt.Sprintf("%v: %v ;", key, value)
	}
	fields["ResponseHeader"] = header
	header = ""
	for key, value := range m{
		header += fmt.Sprintf("%v: %v ;", key, value)
	}
	fields["ResponseBody"] = header
	fields["RequestApi"] = server.Api
	fields["RequestMethod"] = server.Method
	header = ""
	for key, value := range server.Header{
		header += fmt.Sprintf("%v: %v ;", key, value)
	}
	fields["RequestHeader"] = header
	header = ""
	for key, value := range server.Form{
		header += fmt.Sprintf("%v: %v ;", key, value)
	}
	fields["RequestBody"] = header
	fields["RequestTime"] = strconv.FormatFloat(cost, 'g', -1, 64) + "s"
	return fields
}
