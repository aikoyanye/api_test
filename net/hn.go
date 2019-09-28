package net

import (
	"../log"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

/*
	Api：测试的api
	Header：访问api时的请求头，必须是json格式
	Count：并发数
	Method：访问方法
	Form：访问传输的数据，必须是json格式
	Expected：预期返回结果的结构
	DataFormat：传入的参数写在form or url
*/
type Server struct {
	Api 		string					`json:"Api"`
	Header 		map[string]string		`json:"Header"`
	Count 		string					`json:"Count"`
	Method 		string					`json:"Method"`
	Form 		map[string]string		`json:"Form"`
	Expected 	map[string]interface{}	`json:"Expected"`
	DataFormat 	string					`json:"DataFormat"`
}

type SqlInfo struct {
	Type 	string		`json:"Type"`
	Info 	string		`json:"Info"`
	Table 	string		`json:"Table"`
	Fields 	[]string	`json:"Fields"`
}

type CsvInfo struct {
	SavePath 	string		`json:"SavePath"`
	Fields 		[]string	`json:"Fields"`
}

type UploadInfo struct {
	Api 		string					`json:"Api"`
	Header 		map[string]string		`json:"Header"`
	Method 		string					`json:"Method"`
	Fields 		[]string				`json:"Fields"`
}

type ServersLice struct {
	Servers []Server	`json:"Servers"`
	Sql 	SqlInfo		`json:"Sql"`
	Csv 	CsvInfo		`json:"Csv"`
	Upload	UploadInfo	`json:"Upload"`
}

func HttpGo(server *Server) (map[string]interface{}, *http.Response, float64, bool) {
	bodyBuf := &bytes.Buffer{}
	bw := multipart.NewWriter(bodyBuf)
	form := url.Values{}
	var req *http.Request
	var err error

	putData2Form(bw, &form, server.Form)

	content := bw.FormDataContentType()
	bw.Close()
	log.Log.Printf("正在访问API：%v\n", server.Api)
	fmt.Printf("[%v]正在访问API：%v\n", time.Now().Format("2006-01-02 15:04:05"), server.Api)

	// 这步计时
	t := time.Now()
	req, err = http.NewRequest(strings.ToUpper(server.Method), server.Api, nil)
	if err != nil{
		log.Log.Printf("请求失败：%v\n", err)
		fmt.Printf("[%v]请求失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, nil, 0.0, false
	}

	for _, key := range strings.Split(server.DataFormat, "&"){
		if strings.ToLower(key) == "form"{
			req.Body = ioutil.NopCloser(bodyBuf)
		}else if strings.ToLower(key) == "query"{
			req.URL.RawQuery = form.Encode()
		}
	}

	req.Header.Add("Content-Type", content)

	for key, value := range server.Header{
		req.Header.Add(key, value)
	}
	resp, err := http.DefaultClient.Do(req)
	cost := time.Since(t).Seconds()

	if err != nil{
		log.Log.Printf("请求失败：%v\n", err)
		fmt.Printf("[%v]请求失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, nil, 0.0, false
	}
	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil{
		log.Log.Printf("访问API：%v 所返回的Body解析错误：%v\n", server.Api, err)
		fmt.Printf("[%v]访问API：%v 所返回的Body解析错误：%v\n",
			time.Now().Format("2006-01-02 15:04:05"), server.Api, err)
		return nil, nil, 0.0, false
	}
	log.Log.Printf("访问API：%v 成功\n", server.Api)
	fmt.Printf("[%v]访问API：%v 成功\n", time.Now().Format("2006-01-02 15:04:05"), server.Api)
	defer resp.Body.Close()
	return result, resp, cost, true
}

// 将数据写入form中
func putData2Form(bw *multipart.Writer, form *url.Values, fields map[string]string){
	for key, value := range fields{
		if value[0:5] == "file:"{
			fmt.Println(value[5:])
			list := strings.Split(value[5:], "\\")
			if len(list) == 1{
				list = strings.Split(value[5:], "/")
			}
			w, _ := bw.CreateFormFile(key, list[len(list)-1])
			if s, err := os.Open(value[5:]); err == nil{
				io.Copy(w, s)
				s.Close()
			}else{
				log.Log.Printf("将文件：%v写入Form中失败：%v\n", key, err)
				fmt.Printf("[%v]将文件：%v写入Form中失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), key, err)
			}
		}else{
			form.Add(key, value)
			bw.WriteField(key, value)
		}
	}
}