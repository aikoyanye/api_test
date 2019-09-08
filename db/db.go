package info

import (
	"../log"
	"../net"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-mysql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"time"
)

var DBObject *sql.DB

func InitSql(sqlInfo *net.SqlInfo){
	if sqlInfo.Type != "" && sqlInfo.Info != "" && sqlInfo.Table != ""{
		var err error
		DBObject, err = sql.Open(sqlInfo.Type, sqlInfo.Info)
		if err != nil{
			log.Log.Printf("数据库连接错误，可能是连接信息出错：%v\n", err)
			fmt.Printf("[%v]数据库连接错误，可能是连接信息出错：%v\n", time.Now().Format("2006-01-02 15:04:05"), err)
			sqlInfo.Type = ""
			sqlInfo.Info = ""
			sqlInfo.Table = ""
			return
		}
	}
}

func WriteSql(sqlInfo *net.SqlInfo, fields []map[string]string){
	if sqlInfo.Type != "" && sqlInfo.Info != "" && sqlInfo.Table != ""{
		InitSql(sqlInfo)
		log.Log.Println("正在将结果写入数据库...")
		fmt.Printf("[%v]正在将结果写入数据库...\n", time.Now().Format("2006-01-02 15:04:05"))
		for i, fields := range fields{
			values := ""
			for _, field := range sqlInfo.Fields{
				values += "\"" + fields[field] + "\","
			}
			sqlStr := "INSERT INTO " + sqlInfo.Table + " (" + strings.Join(sqlInfo.Fields, ",") + ") VALUES " +
				"(" + values[:len(values)-1] + ")"
			stmt, err := DBObject.Prepare(sqlStr)
			if err != nil{
				log.Log.Printf("写入第 %v 条数据到数据库失败：%v\n", i, err)
				fmt.Printf("[%v]写入第 %v 条数据到数据库失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), i+1, err)
				continue
			}
			res, err := stmt.Exec()
			if err != nil{
				log.Log.Printf("写入第 %v 条数据到数据库失败：%v\n", i, err)
				fmt.Printf("[%v]写入第 %v 条数据到数据库失败：%v\n", time.Now().Format("2006-01-02 15:04:05"), i+1, err)
				continue
			}
			id, _ := res.LastInsertId()
			log.Log.Printf("写入第 %v 条数据到数据库成功：%v\n", i, id)
			fmt.Printf("[%v]写入第 %v 条数据到数据库成功：%v\n", time.Now().Format("2006-01-02 15:04:05"), i+1, id)
		}
	}
}

func CloseDB(sqlInfo *net.SqlInfo){
	if sqlInfo.Type != "" && sqlInfo.Info != "" && sqlInfo.Table != ""{
		DBObject.Close()
	}
}