package main

import (
	"Merryweb/API"
	db "Merryweb/Mysql"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type usolist struct {
	kid string
	vpn []string
}



//获取KIE
func (t *usolist) usokid() usolist {

/*	//uso谷歌验证
	type ReCaptcha struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
		Data    struct {
			TaskId string `json:"taskId"`
		} `json:"data"`
	}


	type taskId struct {
		Success bool   `json:"success"`
		Msg     string `json:"msg"`
		Data    struct {
			Status      string    `json:"status"`
			CaptchaType string    `json:"captchaType"`
			Created     time.Time `json:"created"`
			Response    string    `json:"response"`
		} `json:"data"`
	}



	var kid usolist


	//创建验证码表单
	var ReCapcrete ReCaptcha
	get, err := http.Get("https://api.recaptcha.press/task/create?siteKey=6LcU4RQcAAAAAN7FpoahYQuO9wmoltIJi7O9va9v&siteReferer=https://v3.ujuso.com/&type=ALL&time=ALL&authorization=c45cf09b-75fa-4bee-92e4-a2c0b9ad8359")
	if err != nil {
		log.Println("表单创建失败",err)
		return kid
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		return kid
	}
	err = json.Unmarshal(all,&ReCapcrete)
	if err != nil {
		return usolist{}
	}



	//查询表单验证是否成功
	var Recaptaskid taskId
	for true {

		time.Sleep(1*time.Second)

		get, err := http.Get("https://api.recaptcha.press/task/status?taskId=" + ReCapcrete.Data.TaskId)
		if err != nil {
			log.Println("查询失败:",err)
			return kid
		}

		all, err := io.ReadAll(get.Body)
		if err != nil {
			return kid
		}

		err = json.Unmarshal(all, &Recaptaskid)
		if err != nil {
			return kid
		}


		if Recaptaskid.Data.Status == "Fail" {
			log.Println("查询失败:",Recaptaskid.Msg)
			return kid
		}


		if Recaptaskid.Data.Response != "" {

			key:=struct {
				Code int    `json:"code"`
				Msg  string `json:"msg"`
				Data string `json:"data"`
			}{}


			Resp:=`{"resp":"`+Recaptaskid.Data.Response+`"}`
			post, err := http.Post("https://api.jujuso.com/v1/verify/recaptcha/v2","application/json",strings.NewReader(Resp))
			if err != nil {
				return usolist{}
			}

			readAll, err := io.ReadAll(post.Body)
			if err != nil {
				return usolist{}
			}

			log.Println(string(readAll))


			json.Unmarshal(readAll,&key)
			kid.kid = key.Data

			break
		}

	}*/

	var kid usolist
	kid.kid = "MfoYtJTmhFaGbnSVJNTTBLssejIBfEsc"

	return kid

}

//获取代理VPN
func (t *usolist) https() []string {
/*	//加载HTTPS代理
	_, https := Api.Socket5()
	//分割HTTPS代理
	VPs := strings.Split(https, "\n")
	t.vpn = VPs*/

	VPs:=Api.Freeproxy()
	fmt.Println(VPs)
	return VPs
}

//使用代理访问网页
func (t *usolist) get(kid string,vpn []string) []byte {


	//随机取代理
	dd:=rand.Intn(len(vpn))

	//解析URL HTTPS代理
	parse, err := url.Parse("http://"+vpn[dd])
	if err != nil {
		log.Println(err)
		return nil
	}
	//设置HTTPS客户端
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(parse),
		},
	}

	//获取搜索词库
	phpapi, err := http.Get("http://api.ttos.ltd/api/yiyan/api.php")
	if err != nil {
		return nil
	}

	apis,_:=io.ReadAll(phpapi.Body)

	apisbody:=url.PathEscape(string(apis))

	//访问网页
	get, err := client.Get("https://api.jujuso.com/v1/api_search/vite?page=1&size=1000&kid=" + kid + "&time=ALL&type=ALL&kw="+apisbody)
	if err != nil {
		//log.Println(err)
		return nil
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		return nil
	}


	return all

}

//启动搜索
func usostart()  {
	uso := usolist{

	}

	for {

		vpn := uso.https()
		var kid usolist
		for true {
			kid = uso.usokid()
			if kid.kid != "" {
				break
			}
		}

		sql := db.MysqlMk()

		var wg sync.WaitGroup
		for ss := 0; ss < 64; ss++ {

			wg.Add(1)

			go func() {
				for true {
					//uso 数据结构
					usodoby := struct {
						Code int    `json:"code"`
						Msg  string `json:"msg"`
						Data struct {
							Data []struct {
								CreateTime time.Time `json:"create_time"`
								DiskId     string    `json:"disk_id"`
								DiskName   string    `json:"disk_name"`
								DiskType   string    `json:"disk_type"`
								DiskUser   string    `json:"disk_user"`
								DocId      string    `json:"doc_id"`
								OtherInfo  string    `json:"other_info"`
								SharedTime time.Time `json:"shared_time"`
								UpdateTime time.Time `json:"update_time"`
							} `json:"data"`
							Total int `json:"total"`
						} `json:"data"`
					}{}

					gg := uso.get(kid.kid, vpn)

					json.Unmarshal(gg, &usodoby)

					if usodoby.Msg == "请求错误" {
						log.Println(usodoby.Msg)
						break
					}

					for _, k := range usodoby.Data.Data {

						exec, err := sql.Exec("INSERT INTO `user`.`Doc_id`() VALUES ('" + k.DocId + "')")
						if err != nil {
							log.Println("数据库插入失败", k.DocId)
						}
						if exec != nil {
							log.Println("数据库插入成功", k.DocId)
						}
					}
				}
				wg.Done()
			}()

		}

		wg.Wait()
	}

}


type Uso struct {
	Ld string
	ID int
	DocId string
}


func (Uso *Uso) name(SQL *sql.DB) string  {

	ID, _ := SQL.Query("SELECT * FROM `user`.`JL_` WHERE `UsoID` = 'Uso' LIMIT 0, 1")
	if ID.Next(){
		ID.Scan(&Uso.Ld, &Uso.ID)
	}
	ID.Close()

	SK, _ :=SQL.Query("SELECT * FROM `user`.`Doc_id` WHERE `id` = '"+ strconv.Itoa(Uso.ID)+"'")
	if SK.Next() {
		var ss int
		SK.Scan(&Uso.DocId,&ss)
		SQL.Exec("UPDATE `user`.`JL_` SET `ID` = "+strconv.Itoa(Uso.ID+1)+" WHERE `UsoID` = 'Uso'")
	}
	SK.Close()

	return Uso.DocId


/*	get, err := http.Get("https://api.jujuso.com/v1/disk/info?recaptcha=10&doc_id=")
	if err != nil {
		return
	}

	all, err := io.ReadAll(get.Body)
	if err != nil {
		return
	}

	 sso:=struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			DocId      string `json:"doc_id"`
			DiskId     string `json:"disk_id"`
			DiskType   string `json:"disk_type"`
			DiskPass   string `json:"disk_pass"`
			DiskUser   string `json:"disk_user"`
			DiskName   string `json:"disk_name"`
			OtherInfo  string `json:"other_info"`
			SharedTime int    `json:"shared_time"`
		} `json:"data"`
	}{}

	json.Unmarshal(all,&sso)*/

}

func usego()  {

	var wg sync.WaitGroup
	for i:=0;i<10;i++{
		wg.Add(1)
		go func() {
			SQL:=db.MysqlMk()
			for true{

				JG:= struct {
					Code int    `json:"code"`
					Msg  string `json:"msg"`
					Data struct {
						DocId      string `json:"doc_id"`
						DiskId     string `json:"disk_id"`
						DiskType   string `json:"disk_type"`
						DiskPass   string `json:"disk_pass"`
						DiskUser   string `json:"disk_user"`
						DiskName   string `json:"disk_name"`
						OtherInfo  string `json:"other_info"`
						SharedTime int    `json:"shared_time"`
					} `json:"data"`
				}{}


				D:=Uso{}



				ID:=D.name(SQL)
				var all []byte
				get, err := http.Get("https://api.jujuso.com/v1/disk/info?recaptcha=JJXSjtQzoGrOHrUKJXEKXgfZBkXukbRX&doc_id=" + ID)
				if err != nil {
				}else {
					all, _ = io.ReadAll(get.Body)
					err := json.Unmarshal(all, &JG)
					if err != nil {

					}else {

						exec, _ := SQL.Exec("INSERT INTO `ND_Ujsps` (`Name`, `Url`, `Pass`, `Type`, `User`, `DiskName`, `Time`, `Docid`) VALUES ('" +
							JG.Data.DiskName +
							"', '" +
							JG.Data.DiskId +
							"', '" +
							JG.Data.DiskPass +
							"', '" +
							JG.Data.DiskType +
							"', '" +
							JG.Data.DiskUser +
							"', '" +
							JG.Data.OtherInfo +
							"', '" +
							time.Unix(int64(JG.Data.SharedTime), 0).Format("2006-01-02 15:04:05") +
							"', '" +
							ID +
							"')")

						if exec != nil {
							log.Println("网盘数据插入成功:",ID,exec)
						}
					}
				}



			}

			wg.Done()
		}()
	}
	wg.Wait()
}

func main() {


	go usostart()
	go usego()

	time.Sleep(1000000000*time.Second)


	// go usostart()
}








