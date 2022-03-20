package main

import (
	"Merryweb/API"
	db "Merryweb/Mysql"
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"

)

func main() {

	//加载静态资源文件
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("./Web/resource/"))))

	//访问首页登陆页面
	http.HandleFunc("/", login)
	//访问用户注册页面
	http.HandleFunc("/register", register)
	//访问工具箱页面
	http.HandleFunc("/tools", tools)

	//加载API项目
	Api.Server()

	//建立web
	err := http.ListenAndServe(":25565", nil)
	if err != nil {
		log.Println(err)
	}
}

//访问首页登陆
func login(w http.ResponseWriter, r *http.Request) {

	b, p := db.Mysqllcookiecx(r)
	log.Println(p)
	//Cookie登陆
	if b {
		t, _ := template.ParseFiles("./Web/tools.html")
		t.Execute(w, "登陆成功! 用户名:"+p.User)
		log.Println("远程终端", r.RemoteAddr, "Cookie登陆:", r.Header["Cookie"])
		return
	}

	//解析表单
	r.ParseForm()

	//判断是否登陆
	if r.Form["user"] == nil && r.Form["pass"] == nil {
		t, _ := template.ParseFiles("./Web/login.html")
		t.Execute(w, "请输入用户和密码!")
		return
	} else {
		if db.Mysqlu(r) {

			//MD5加密COOKE
			Cookepassmd5 := md5.Sum([]byte(r.Form["pass"][0]))

			//设置COOKE
			Cooke := http.Cookie{
				Name:  r.Form["user"][0],
				Value: hex.EncodeToString(Cookepassmd5[:16]),
			}
			http.SetCookie(w, &Cooke)

			//向数据库添加cookie
			db.Mysqlcookieadd(r.Form["user"][0], r.Form["user"][0]+"="+hex.EncodeToString(Cookepassmd5[:16]))

			//登陆成功执行操作
			t, _ := template.ParseFiles("./Web/tools.html")
			t.Execute(w, "登陆成功! 用户名:"+r.Form["user"][0])
			log.Println("远程终端", r.RemoteAddr, "正在登陆用户", "用户名:", r.Form["user"], "用户密码:", r.Form["pass"])
			return
		} else {
			//登陆失败执行操作
			t, _ := template.ParseFiles("./Web/login.html")
			t.Execute(w, "用户名或密码错误!")
			return
		}
	}

}

//访问注册网站
func register(w http.ResponseWriter, r *http.Request) {
	//解析表单数据
	r.ParseForm()
	//获取用户数据
	user := r.Form["user"]
	pass := r.Form["pass"]

	//判断是否能注册
	if user == nil && pass == nil {
		t, _ := template.ParseFiles("./Web/register.html")
		t.Execute(w, "请输入注册用户名邮箱!")
		return
	} else if len(pass) == 1 {
		t, _ := template.ParseFiles("./Web/register.html")
		t.Execute(w, "请输入确认密码!")
		return
	} else if pass[0] != pass[1] {
		t, _ := template.ParseFiles("./Web/register.html")
		t.Execute(w, "密码不一致!")
		return
	} else if pass[0] == pass[1] {
		//db.Mysqluseradd(r)
		//判断数据库是否插入
		if true {
			t, _ := template.ParseFiles("./Web/login.html")
			t.Execute(w, "暂时不给予注册！")
			return
		} else {
			t, _ := template.ParseFiles("./Web/register.html")
			t.Execute(w, "注册失败用户名,已经存在!")
			return
		}

	} else {
		t, _ := template.ParseFiles("./Web/register.html")
		t.Execute(w, "内部数据校验错误!")
		return
	}

}

//访问工具箱页面
func tools(w http.ResponseWriter, r *http.Request) {

	//解析表单
	r.ParseForm()
	b, p := db.Mysqllcookiecx(r)

	//Cookie登陆
	if b {
		t, _ := template.ParseFiles("./Web/tools.html")
		t.Execute(w, "登陆成功! 用户名:"+p.User)
		log.Println("远程终端", r.RemoteAddr, "Cookie登陆:", r.Header["Cookie"])
		return
	} else {
		t, _ := template.ParseFiles("./Web/login.html")
		t.Execute(w, "")
		return
	}

}
