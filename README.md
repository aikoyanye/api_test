 # 非常简单的、估计没啥用处的api自动化测试工具
 能够自定义访问Request的参数（method、header、form），获取访问的时间，并对比返回的结果是不是符合预期（如果需要的话）

## 使用方法
go build main.go 
-path [json案例文件路径]   
-mode [运行模式，默认：default]   
-log [保存的log文件路径，默认：程序当前路径]   
[xx] hh [xxx] hhh ...（可变长参数，将json案例中的字符串替换成其他字符串，[xx]替换成hh，不限格式，不能以“-”开头）

### -mode [default][parallel][performance][all]
default：在主线程中，按照顺序与次数，访问API  
parallel：并行的访问多个API（非相同API）  
performance：高并发的访问相同的API（可能有问题）  
all：并行且高并发的访问每个API（同样可能有问题）  

### 案例样式，只接受JSON数据，且字段要符合相同（包括大小写），在根目录的1.json就是例子
Servers：是个列表，需要访问的API按照顺序排列即可

以下是每组API的说明：  
Api：待访问的url  
Header：访问url时的请求头，不需要就留空"Header": {}，必须是string: string  
Count：访问url的次数，性能模式时为并发数
Method：访问url的方法  
Form：访问url时传递的表单数据，不需要就留空"Form": {}，也必须是string: string  
Expected：期望返回的结果，key要准确，value随意，string: string or string: list or string: map  

### 将结果写入CSV或SQL中  
支持写入的字段信息："RequestBodyExpected", "RequestTime", "RequestBody", "RequestHeader", "RequestMethod", "RequestApi", "ResponseBody",
      "ResponseHeader", "ResponseProto", "ResponseStatusCode", "ResponseStatus", "ResponseContentLength", "ResponseTransferEncoding"  
写入CSV：  
SavePath：保存Csv文件的路径，要精确到文件名和后缀  
Fields：需要写入的字段（可能不是全部都需要吧）  
写入Sql：  
Type：目前只支持mysql、sqlite3  
Info：连接sql的信息，例如：root:rootroot@tcp(127.0.0.1:3306)/test?charset=utf8（mysql） or log.db（sqlite3）  
Table：待插入的数据库的表名（必须有字段，不会自动帮你创建）  
Fields：需要写入的字段（可能不是全部都需要吧）  
