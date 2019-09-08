package net

import (
	"../log"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
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
*/
type Server struct {
	Api 		string					`json:"Api"`
	Header 		map[string]string		`json:"Header"`
	Count 		int						`json:"Count"`
	Method 		string					`json:"Method"`
	Form 		map[string]string		`json:"Form"`
	Expected 	map[string]interface{}	`json:"Expected"`
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

type ServersLice struct {
	Servers []Server	`json:"Servers"`
	Sql 	SqlInfo		`json:"Sql"`
	Csv 	CsvInfo		`json:"Csv"`
}

func HttpGo(server *Server) (map[string]interface{}, *http.Response, float64, bool) {
	bodyBuf := &bytes.Buffer{}
	bw := multipart.NewWriter(bodyBuf)
	form := url.Values{}
	var req *http.Request
	var err error

	for key, value := range server.Form{
		form.Add(key, value)
		bw.WriteField(key, value)
	}
	content := bw.FormDataContentType()
	bw.Close()
	log.Log.Printf("正在访问API：%v\n", server.Api)
	fmt.Printf("[%v]正在访问API：%v\n", time.Now().Format("2006-01-02 15:04:05"), server.Api)

	t := time.Now()
	// 这步计时
	req, err = http.NewRequest(strings.ToUpper(server.Method), server.Api, nil)
	if err != nil{
		log.Log.Printf("请求失败：%v\n", err)
		fmt.Printf("[%v]请求失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, nil, 0.0, false
	}
	req.URL.RawQuery = form.Encode()
	req.Body = ioutil.NopCloser(bodyBuf)
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
		log.Log.Printf("访问API：%v所返回的Body解析错误：%v\n", server.Api, err)
		fmt.Printf("[%v]访问API：%v所返回的Body解析错误：%v\n",
			time.Now().Format("2006-01-02 15:04:05"), server.Api, err)
		return nil, nil, 0.0, false
	}
	log.Log.Printf("访问API：%v 成功\n", server.Api)
	fmt.Printf("[%v]访问API：%v 成功\n", time.Now().Format("2006-01-02 15:04:05"), server.Api)
	defer resp.Body.Close()
	return result, resp, cost, true
}
