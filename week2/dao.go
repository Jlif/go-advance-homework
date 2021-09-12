package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	DBHostsIp  = "localhost:3306"
	DBUserName = "root"
	DBPassWord = "123456"
	DBName     = "test"
)

func main() {
	//连接至数据库
	db, err := sql.Open("mysql", DBUserName+":"+DBPassWord+"@tcp("+DBHostsIp+")/"+DBName)
	if err != nil {
		fmt.Println(err)
		panic("连接数据库失败")
	}
	err = insert(db)
	if err != nil {
		fmt.Println(err)
	}

	err = query(db)
	if err != nil {
		fmt.Println(err)
	}
	//关闭数据库连接
	err = db.Close()
	if err != nil {
		fmt.Println(err)
	}
}

//插入demo
func insert(db *sql.DB) error {
	//准备插入操作
	stmt, err := db.Prepare("INSERT demo (id,name) values (?,?)")
	if err != nil {
		return errors.Wrap(err, "数据库准备插入操作异常")
	}
	//执行插入操作
	res, err := stmt.Exec(1, "jlif")
	if err != nil {
		return errors.Wrap(err, "数据库插入操作异常")
	}
	//返回最近的自增主键id
	id, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "返回最近的自增主键id异常")
	}
	fmt.Println("LastInsertId: ", id)
	return nil
}

// query 查询demo
func query(db *sql.DB) error {
	//rows：返回查询操作的结果集
	rows, err := db.Query("SELECT * FROM demo")
	if err != nil {
		return errors.Wrap(err, "返回查询操作的结果集异常")
	}
	//第一步：接收在数据库表查询到的字段名，返回的是一个string数组切片
	columns, _ := rows.Columns() // columns:  [user_id user_name user_age user_sex]
	//根据string数组切片的长度构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		//将查询到的字段名的地址复制到scanArgs数组中
		err = rows.Scan(scanArgs...)
		if err != nil {
			return errors.Wrap(err, "遍历结果集异常")
		}
		//将行数据保存到record字典
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				//字段名 = 字段信息
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
	return nil
}
