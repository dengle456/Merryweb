package Api

import (
	db "Merryweb/Mysql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type name struct {
	Code int
	Msg  string
	Data []struct{
		Name string
		Url  string
		Pass string
		Tpye string
		USER string
		Diskname string
		Size string
		TIME string
		UUid string
	}
}



//tools 云盘api
func bdyapi(w http.ResponseWriter,r *http.Request)  {

	d:=name{}
	r.ParseForm()



	switch {
	case r.Form["name"] != nil && r.Form["size"] != nil && r.Form["page"] != nil:
		page, _ := strconv.Atoi(r.Form["page"][0])
		size, _ := strconv.Atoi(r.Form["size"][0])
		r.Form["page"][0] = strconv.Itoa(page * size - size)
		//S:="SELECT * FROM `user`.`ND_Ujsps` WHERE `Name` LIKE '%"+ r.Form["name"][0] + "%' OR `DiskName` LIKE '%"+r.Form["name"][0] + "%' LIMIT "+r.Form["page"][0]+", "+r.Form["size"][0]+""
		G:= d.Json("SELECT * FROM `user`.`ND_Ujsps` WHERE `Name` LIKE '%"+r.Form["name"][0]+"%' OR `DiskName` LIKE '%"+r.Form["name"][0]+"%' LIMIT "+r.Form["page"][0]+", "+r.Form["size"][0]+"")

		fmt.Fprintf(w, string(G))
		return
	}
	fmt.Fprintf(w,`{"Code":301,"Msg":"无请求内容","Data":null}`)
}

func (dode *name)Json(sql string)[]byte  {

	SQL:=db.MysqlMk()


	Q,err:=SQL.Query(sql)
	if err != nil {
		log.Println(err)
		dode.Code = 404
		dode.Msg = "请求参数错误"
		dode.Data =nil
		marshal, _ := json.Marshal(dode)
		return marshal
	}
	var n1,n2,n3,n4,n5,n6,n7,n8,n9 string
	for i:=0;Q.Next();i++{
		dode.Data = append(dode.Data,make([]struct{
			Name string
			Url  string
			Pass string
			Tpye string
			USER string
			Diskname string
			Size string
			TIME string
			UUid string
		},1)...)

		Q.Scan(&n1,&n2,&n3,&n4,&n5,&n6,&n7,&n8,&n9)
		dode.Data[i].Name = n1
		dode.Data[i].Pass = n3
		log.Println(n7)
		switch  n4{
		case "BDY":
			dode.Data[i].Tpye =  "百度云盘"
			dode.Data[i].Url = "https://pan.baidu.com/share/init?surl=" + n2
		case "ALY":
			dode.Data[i].Tpye =  "阿里云盘"
			dode.Data[i].Url = "https://www.aliyundrive.com/s/" + n2
		case "LZY":
			dode.Data[i].Tpye =  "蓝奏云盘"
			dode.Data[i].Url = "https://lanzoux.com/" + n2
		case "QUARK":
			dode.Data[i].Tpye =  "夸克云盘"
			dode.Data[i].Url = "https://pan.quark.cn/s/" + n2
		}
		dode.Data[i].USER = n5
		dode.Data[i].Diskname = n6
		dode.Data[i].Size = n7
		dode.Data[i].TIME = n8
		dode.Data[i].UUid = n9
	}
	Q.Close()

	if dode.Data != nil {
		dode.Code = 200
		dode.Msg = "请求成功"
		marshal, err := json.Marshal(dode)
		if err != nil {
			fmt.Println(err)
		}
		return marshal
	}else {
		dode.Code = 303
		dode.Msg = "无搜索结果"
		marshal, err := json.Marshal(dode)
		if err != nil {
			fmt.Println(err)
		}
		return marshal
	}
}
