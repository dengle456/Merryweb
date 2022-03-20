package Api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Proxy []struct {
	Id             int       `json:"id"`
	LocalId        int       `json:"local_id"`
	ReportId       string    `json:"report_id"`
	Addr           string    `json:"addr"`
	Type           int       `json:"type"` //TYPE 1.http 2.https 3.socks4 4.socks5
	Kind           int       `json:"kind"`
	Timeout        int       `json:"timeout"`
	Cookie         bool      `json:"cookie"`
	Referer        bool      `json:"referer"`
	Post           bool      `json:"post"`
	Ip             string    `json:"ip"`
	AddrGeoIso     string    `json:"addr_geo_iso"`
	AddrGeoCountry string    `json:"addr_geo_country"`
	AddrGeoCity    string    `json:"addr_geo_city"`
	IpGeoIso       string    `json:"ip_geo_iso"`
	IpGeoCountry   string    `json:"ip_geo_country"`
	IpGeoCity      string    `json:"ip_geo_city"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Skip           bool      `json:"skip"`
	FromCache      bool      `json:"from_cache"`
}

var ip Proxy
var ferr []string


//获取checkerproxy.net 到JSON文本
func jsoncache() {
	//获取当前时间到文本
	Protime := time.Now().Truncate(time.Hour * 12).String()
	log.Println(Protime)

	ProxyStr := strings.Split(Protime, " ")
	ProxyStrget := strings.ReplaceAll(ProxyStr[0], "/", "-")

	//获取代理
	dd, _ := http.Get("https://checkerproxy.net/api/archive/" + ProxyStrget)
	all, _ := ioutil.ReadAll(dd.Body)
	json.Unmarshal(all, &ip)
	if &ip != nil {
		log.Println("代理JSON获取成功!")
	}else {
		log.Println("JSON文本获取失败")
	}
}

//获取free-proxy.cz	到字符串
func FerrSocks() []string {

	//获取页码
	ss:= Getsing("http://free-proxy.cz/zh/proxylist/country/all/socks5/ping/all")
	Azeem := regexp.MustCompile(`<a href="/zh/proxylist/country/all/socks5/ping/all/.*?</a>`)
	ym:=len(Azeem.FindAllString(ss,-1))

	//获取网页数据
	for y:=2;y<ym;y++ {
		ss = ss + Getsing("http://free-proxy.cz/zh/proxylist/country/all/socks5/ping/all/" +strconv.Itoa(y))
	}

	//匹配ip数据
	str1:= regexp.MustCompile(`Base64.decode\(".*?"\)`)
	fat1:=str1.FindAllString(ss,-1)
	//匹配端口数据
	str2:=regexp.MustCompile(`style=''>.*?</span>`)
	fat2:=str2.FindAllString(ss,-1)

	//获取文本所有IP长度
	cd:=len(fat1)
	cd1:=len(fat2)
	if cd == cd1 {
		//组合IP和端口
		for i:=0;i<cd;i++ {
			//将所有IP解码
			fat1[i]=strings.ReplaceAll(fat1[i],"Base64.decode(\"","")
			fat1[i]=strings.ReplaceAll(fat1[i],"\")","")
			ff, _ :=base64.StdEncoding.DecodeString(fat1[i])
			jm:=string(ff)


			//将所有端口筛选
			fat2[i]=strings.ReplaceAll(fat2[i],"style=''>","")
			fat2[i]=strings.ReplaceAll(fat2[i],"</span>","")


			//将所有IP和端口组合
			fat1[i] = jm + ":" + fat2[i]
		}
		log.Println("ferrpoxy获取成功")
		return fat1

	}else {
		log.Println("数据不统一")
	}

	return nil


}

//获取IP freeproxy.world
func Freeproxy() []string {

	var f []byte
	var str string
	var all string

	for i := 0; i < 21; i++ {
		d, _ := http.Get("https://www.freeproxy.world/?type=http&anonymity=&country=&speed=&port=&page=" + strconv.Itoa(i))
		if d != nil {
			log.Println("代理网页访问成功")
			f, _ = ioutil.ReadAll(d.Body)
			all = string(f)
			str = str + all

		} else {
			log.Println("代理网页访问失败")
		}

	}

	ze := regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	fz := ze.FindAllString(str, -1)

	ff := "port=.*\""
	ze = regexp.MustCompile(ff)
	dk := ze.FindAllString(str, -1)
	log.Println("当前代理长度:", len(dk))

	s := make([]string, len(dk))

	var dd string
	for i := 0; i < len(dk); i++ {

		ffz := fz[i]
		ffe := dk[i]

		ffe = strings.Replace(ffe, "port=", "", -1)
		ffe = strings.Replace(ffe, "\"", "", -1)

		s[i] = ffz + ":" + ffe
		//log.Println("成功获取IP", s[i])
		dd = dd + s[i] + "\n"

	}

	return s

}

//GET URI返回字符串
func Getsing(uri string) string {
	Client:=http.Client{}
	res,_:=http.NewRequest("GET",uri,nil)
	res.Header.Add("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.82 Safari/537.36")
	res.Header.Add("Content-Type","application/json;charset=UTF-8")
	dd,_:=Client.Do(res)
	ff,_:=ioutil.ReadAll(dd.Body)
	defer dd.Body.Close()
	return string(ff)

}

//获取socks5 和 HTTPS 代理
func Socket5() (string, string) {

	jsoncache()

	//筛选可用代理
	var wg sync.WaitGroup
	var Sock5 string
	var https string

	//获取checkerproxy.net 可用代理
	for i := 0; len(ip) > i; i++ {
		i := i
		wg.Add(1)
		//获取socks5
		go func() {
			if ip[i].Type == 4 {
				dialer := net.Dialer{
					Timeout: 5000 * time.Millisecond,
				}
				Cliet, err := dialer.Dial("tcp", ip[i].Addr)
				if err != nil {
					log.Println("建立连接失败:", ip[i].Addr)
				} else {
					Cliet.SetReadDeadline(time.Now().Add(1 * time.Second))
					Cliet.Write([]byte{0x05, 0x01, 0x00})
					Wrbyte := make([]byte, 1024)
					read, _ := Cliet.Read(Wrbyte)
					log.Printf("%v %v %v % x", "建立连接成功:", Cliet.RemoteAddr().String(), "返回字节:", Wrbyte[:read])
					if Wrbyte[0] == 0x05 && Wrbyte[1] == 0x00 {
						log.Printf("%s%s", Cliet.RemoteAddr().String(), "认证成功")
						if Sock5 != "" {
							Sock5 = Sock5 + ip[i].Addr + "\n"
						} else {
							Sock5 = ip[i].Addr + "\n"
						}
					}
				}
			}
			wg.Done()
		}()

		//获取https
		if ip[i].Type == 2 {
			if https != "" {
				https = https + ip[i].Addr + "\n"
			} else {
				https = ip[i].Addr + "\n"
			}
		}
	}
	//获取free-proxy.cz 可用代理
	for i :=0;len(ferr) > i ;i++{
		wg.Add(1)
		i := i
		go func() {
			dialer := net.Dialer{
				Timeout: 5000 * time.Millisecond,
			}
			Cliet, err := dialer.Dial("tcp", ferr[i])
			if err != nil {
				log.Println("建立连接失败:", ferr[i])
			} else {
				Cliet.SetReadDeadline(time.Now().Add(1 * time.Second))
				Cliet.Write([]byte{0x05, 0x01, 0x00})
				Wrbyte := make([]byte, 1024)
				read, _ := Cliet.Read(Wrbyte)
				log.Printf("%v %v %v % x", "建立连接成功:", Cliet.RemoteAddr().String(), "返回字节:", Wrbyte[:read])
				if Wrbyte[0] == 0x05 && Wrbyte[1] == 0x00 {
					log.Printf("%s%s", Cliet.RemoteAddr().String(), "认证成功")
					if Sock5 != "" {
						Sock5 = Sock5 + ferr[i] + "\n"
					} else {
						Sock5 = ferr[i] + "\n"
					}
				}
			}
			wg.Done()
		}()


	}



	wg.Wait()






	return Sock5, https
}
