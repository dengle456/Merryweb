package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

func main() {


	open:=opensql()

	QS,err:=open.Query("SELECT * FROM `user`.`heimisou`")
	if err != nil {
		fmt.Println(err)
	}

	type SALDADA struct {
		Name     string
		url      string
		pass     string
		time     string
		size     string
		Diskinfo string
		UUID     string
	}


	FD:= SALDADA{}
	for QS.Next(){
		err := QS.Scan(&FD.Name, &FD.url, &FD.pass, &FD.time, &FD.size,&FD.Diskinfo,&FD.UUID)
		if err != nil {
			log.Println(err)
			return
		}


		FD.url=strings.ReplaceAll(FD.url,"https://pan.baidu.com/s/1","")
		FD.UUID = "BDY|"+FD.url


		_, err = open.Exec("INSERT INTO `user`.`ND_Ujsps`(`Name`, `Url`, `Pass`, `Type`, `User`, `DiskName`, `size`, `Time`, `Docid`) " +
			"VALUES ('" + FD.Name + "', '" + FD.url + "', '" + FD.pass + "', 'BDY', 'admin', '" + FD.Diskinfo + "', '" + FD.size + "', '" + FD.time + "', '" + FD.UUID + "')")
		if err != nil {
			fmt.Println(err)
		}else {
			log.Println("数据添加成功",FD.UUID)
		}








	}




}

func opensql()*sql.DB {
	dd,err:=sql.Open("mysql","user:ssD790416@tcp(192.168.0.7)/user")
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(`数据库连接成功!`)
	}
	return dd
}