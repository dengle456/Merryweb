package Api

import (
	"Merryweb/Mysql"
	"encoding/json"
	"github.com/gorilla/websocket"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)


type Ms struct {
	Name string
	Dody string
	Msg string
}


//加载项目
func Server()  {

	//加载SQP
	http.HandleFunc("/Geturl",get)

	//加载SOCKS5
	http.HandleFunc("/Socks5",ApiSocks5)

	//加载Q绑查询
	http.HandleFunc("/Qb",Qb)

	//加载身份证查询
	http.HandleFunc("/sfz",sfz)

	//加载百度云搜
	http.HandleFunc("/bdy",bdy)
	http.HandleFunc("/bdyapi",bdyapi)

	//websocket
	http.HandleFunc("/ws",ws)




}


//websocket
func ws(w http.ResponseWriter,r *http.Request) {




	//http升级websocket
	d:=websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	fd,_:=d.Upgrade(w,r,nil)


	//发送信息告诉客户端连接成功!
	fd.WriteMessage(1,[]byte("1.连接成功,建立成功"))

	//回调信息
	go wsmsg(fd)








}

//websocket信息回调
func wsmsg(w *websocket.Conn)  {

	for true {

		_, i, err := w.ReadMessage()
		if err != nil {
			log.Println(w.RemoteAddr(),err)
			return
		}

		log.Println(w.RemoteAddr(),":" + string(i))


		if strings.Index(string(i),"http") != -1 && strings.Index(string(i),"https") != -1{
			response, err := http.Get(string(i))
			if err != nil {
				return
			}
			all, err := io.ReadAll(response.Body)
			if err != nil {
				return
			}
			w.WriteMessage(1,all)
		}


	}

}

//tools 访问GET
func get(w http.ResponseWriter,r *http.Request)  {


	//解析表单
	r.ParseForm()

	//加载模板
	t, _ := template.ParseFiles("./Web/Geturl.html")

	//查询COOKIE
	d,p:=db.Mysqllcookiecx(r)

	//设置发送数据
	msg :=Ms{
		p.User,
		"",
		"请输入GET地址",
	}


	//如果URL表单有值,GET地址
	if d && r.Form["url"] != nil{
		msg.Dody = "a"
		t.Execute(w,msg)
		return
		}


	//如果没有值
	if d{
		t.Execute(w,msg)
		return
	}


}

//tools 访问Socks5
func ApiSocks5(w http.ResponseWriter,r *http.Request)  {

	//解析表单
	r.ParseForm()

	//加载模板
	t, _ := template.ParseFiles("./Web/Socks5.html")

	//查询COOKIE
	d,p:=db.Mysqllcookiecx(r)

	//设置发送数据
	msg :=Ms{
		p.User,
		"",
		"Socks5",
	}


	//如果URL表单有值,GET地址
	if d && r.Form["Socks5"] != nil{

		if r.Form["Socks5"][0] == "Socks"{
			Dody,_:=Socket5()
			if Dody == "" {
				msg.Dody = time.Now().String() + "暂时无数据!"
			}else {
				msg.Dody = Dody
			}
			t.Execute(w,msg)
			return
		}else if r.Form["Socks5"][0] == "https" {
			_,Dody:=Socket5()
			msg.Dody = Dody
			t.Execute(w,msg)
			return
		}else {
			msg.Dody = "口令输入错误"
			t.Execute(w,msg)
			return
		}


	}

	//如果没有值
	if d{
		t.Execute(w,msg)
		return
	}


}

//tools 访问Q绑查询
func Qb(w http.ResponseWriter,r *http.Request)  {
	//解析表单
	r.ParseForm()

	//加载模板
	t, _ := template.ParseFiles("./Web/Qb.html")

	//查询COOKIE
	d,p:=db.Mysqllcookiecx(r)

	//设置发送数据
	type T struct {
		Name string
		Status    int    `json:"status"`
		Message   string `json:"message"`
		Phone     string `json:"phone"`
		Phonediqu string `json:"phonediqu"`
		Lol       string `json:"lol"`
		Wb        string `json:"wb"`
		Qqlm      string `json:"qqlm"`
	}
	QQ:= T{}

	//如果URL表单有值,GET地址
	if d && r.Form["QQ"] != nil{
		response, err := http.Get("https://zy.xywlapi.cc/qqcx?qq=" + r.Form["QQ"][0])
		if err != nil {
			QQ.Message = err.Error()
			t.Execute(w,QQ)
			return
		}

		all, err := io.ReadAll(response.Body)
		if err != nil {
			QQ.Message = err.Error()
			t.Execute(w,QQ)
			return
		}

		json.Unmarshal(all,&QQ)
		QQ.Name = p.User
		t.Execute(w,QQ)
		return
	}

	//如果没有值
	if d{
		QQ.Name = p.User
		t.Execute(w,QQ)
		return
	}

}

//tools 访问身份证查询
func sfz(w http.ResponseWriter,r *http.Request)  {

	type msg struct {
		Name string
		Xm []string
		Sfz []string
	}
	//解析表单
	r.ParseForm()
	//加载模板
	t, _ := template.ParseFiles("./Web/sfz.html")


	//COOKIE认证
	Cookie,name:=db.Mysqllcookiecx(r)


	if Cookie && r.Form["sfz"] != nil {
		mx:=db.Mysqlsfzcx(r)
		Data:=msg{
			Name: name.User,
			Xm: mx.Xm,
			Sfz:mx.Sfz,
		}
		//log.Println(Data.Xm)
		t.Execute(w,Data)
		return

	}



	//认证通过输出
	if Cookie {
		mx:=db.Mysqlsfzcx(r)
		//输出数据结构
		Data:=msg{
			Name: name.User,
			Xm: mx.Xm,
			Sfz:mx.Sfz,
		}

		t.Execute(w,Data)
		return

	}




}

//tools 访问盘搜百度
func bdy(w http.ResponseWriter,r *http.Request)  {


	r.ParseForm()
	t,user:=db.Mysqllcookiecx(r)
	log.Println(user)

	//解析表当
	if t {
		t,_:=template.ParseFiles("./Web/bdy.html")
		t.Execute(w,user.User)
	}else {
		t,_:=template.ParseFiles("./Web/bdy.html")
		t.Execute(w,"游客")
	}


}






