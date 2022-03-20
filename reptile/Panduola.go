package main

import (
	db "Merryweb/Mysql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Panduolas struct {
	dd int
	ID []struct {
		ID   string
		Surl string
		Name string
		time string
		pass string
	}
}


func main() {



	var ss struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Id       int    `json:"id"`
			Password string `json:"password"`
			Surl     string `json:"surl"`
			Title    string `json:"title"`
			Isdir    string `json:"isdir"`
			Ctime    int    `json:"ctime"`
		} `json:"data"`
	}


	g:=db.MysqlMk()

	for i:=0;i<64;i++ {


		go func() {
			for true {
			Panduola:=Pandolalist()

			for _,K:=range Panduola.ID {

				get, err := http.Get("http://159.75.239.59:7987/api/cxtquery?clienttype=0&referral=&surl=" + K.Surl + "&t=" + K.ID + "&version=0.0.4")
				if err != nil {
					fmt.Println(err)
				}else {
					all, err := io.ReadAll(get.Body)
					if err != nil {
						fmt.Println(err)
						return
					}
					json.Unmarshal(all, &ss)
				}

				K.Surl = strings.TrimLeft(K.Surl,"1")

				if ss.Data != nil {
					K.pass = ss.Data[0].Password
					SQL:="INSERT INTO `user`.`ND_Ujsps`(`Name`, `Url`, `Pass`, `Type`, `User`, `DiskName`, `size`, `Time`, `Docid`) " +
						"VALUES ('"+K.Name+"', '"+K.Surl+"', '"+K.pass+"', 'BDY', 'admin', '"+K.Name+"', '', '"+K.time+"', 'BDY|"+ K.Surl +"')"
					_, err := g.Exec(SQL)
					if err != nil {
						log.Println(err)
					}else {
						log.Println("数据库添加成功：",K.Surl)
					}
				}

			}
			}
		}()

	}

	for i:=0;i<3000;i++ {
		time.Sleep(time.Hour)
		log.Println("已经运行:",i,"分钟")
	}


}

func Pandolalist() Panduolas {

	Panduola:= Panduolas{}

	//search.pandown.me Json数据结构
	Que:=struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Id       int    `json:"id"`
			Shareid  int 	`json:"shareid"`
			Uk       int 	`json:"uk"`
			Password string `json:"password"`
			Surl     string `json:"surl"`
			Title    string `json:"title"`
			Isdir    string `json:"isdir"`
			Ctime    int    `json:"ctime"`
		} `json:"data"`
	}{}

	//随机获取文字API
	response, err := http.Get("http://api.ttos.ltd/api/yiyan/api.php")
	if err != nil {
		return Panduolas{}
	}
	buf,_:=io.ReadAll(response.Body)
	str:=strings.ReplaceAll(string(buf),"\n","")
	str =url.PathEscape(str)


	//获取JSON格式存入变量结构体
	var ss int
	var sd int
	for i:=0;i<20;i++ {
		get, err := http.Get("https://search.pandown.me/api/query?clienttype=0&version=0.0.4&highlight=10&key="+ str +"&page="+ strconv.Itoa(i) +"&timestamp=1647620108&sign=136ba3012f121f2b7d602d876ae86b4e")
		if err != nil {
			return Panduolas{}
		}
		buf,_=io.ReadAll(get.Body)
		err = json.Unmarshal(buf, &Que)
		if err != nil {
			log.Println(err)
			return Panduolas{}
		}

		g:=len(Que.Data)
		sd=len(Que.Data)
		if i == 0 {
			 ss = 0
		}else if i == 1 {
			ss = sd
		}else {
			if g != 0 {
				ss = ss + g
			}
		}


		for v,k:= range Que.Data{
			//log.Println(ss+v)
			Panduola.ID=append(Panduola.ID,make([]struct{
				ID   string
				Surl string
				Name string
				time string
				pass string
			},1)...)

			Panduola.ID[ss+v].ID = strconv.Itoa(k.Id)
			Panduola.ID[ss+v].Surl = k.Surl
			Panduola.ID[ss+v].Name = k.Title
			Panduola.ID[ss+v].time = ctime(k.Ctime)
		}

	}

	return Panduola
}

func ctime(timestamp int) string  {

	timeobj := time.Unix(int64(timestamp), 0)

	date := timeobj.Format("2006-01-02 15:04:05")

	return date
}