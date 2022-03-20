package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)


//Cookie结构
type name struct {
	User string
	Cookie string
}

//身份证结构
type sfz struct {
	Xm []string
	Sfz []string
}



func MysqlMk() *sql.DB  {

	db, err := sql.Open("mysql", "user:ssD790416@tcp(localhost:3306)/user")
	if err != nil {
		log.Println(err)
	}
	_, err = db.Exec("CREATE TABLE `user`.`usercookie`  (\n  `user` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,\n  `cookie` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,\n  PRIMARY KEY (`user`) USING BTREE,\n  UNIQUE INDEX `user`(`cookie`) USING BTREE\n);")
	_, err = db.Exec("CREATE TABLE `user`.`user`  (\n  `id` int(10) NOT NULL AUTO_INCREMENT,\n  `Email` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,\n  `user` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,\n  `pass` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,\n  PRIMARY KEY (`id`, `user`) USING BTREE,\n  UNIQUE INDEX `user`(`user`) USING BTREE\n);")
	if err.Error() == "Error 1050: Table 'user' already exists" {
		//log.Println("SQL表,已存在或已经创建")
	}else {
		log.Println(err)
	}

	return db
}

//向数据库添加用户
func Mysqluseradd(r *http.Request) bool   {
	log.Println("远程终端",r.RemoteAddr,"正在注册用户","用户名:",r.Form["user"],"用户密码:",r.Form["pass"],"邮箱:",r.Form["email"])
	db:=MysqlMk()
	_, err := db.Exec("INSERT INTO `user` (`user`, `pass`, `Email`) VALUES ( '"+r.Form["user"][0]+"', '"+ r.Form["pass"][0]+"', '"+r.Form["email"][0] +"')")
	defer db.Close()
	if err != nil {
		log.Println("远程终端",r.RemoteAddr,"MYSQL错误:",err)
		return false
	}else {
		log.Println("用户注册成功",r.Form["user"],r.Form["email"][0])
		return true
	}


}

//查询数据库用户是否存在
func Mysqlu(from *http.Request)  bool {

	db:= MysqlMk()

	exec, err := db.Query("SELECT * FROM `user`.`user` WHERE `user` LIKE '"+ from.Form["user"][0]+"' AND `pass` LIKE '" +
		from.Form["pass"][0] +
		"' LIMIT 0, 1")
	if err != nil {
		log.Println(err)
		return false
	}

	defer db.Close()
	return exec.Next()

}

//数据库添加COOKIE
func Mysqlcookieadd(user string,cookie string)  {
	db:= MysqlMk()
	_, err := db.Exec("INSERT INTO `user`.`usercookie`(`user`, `cookie`) VALUES ('"+ user +"', '"+ cookie+"')")
	if err != nil {
		log.Println(err)
		return
	}

	defer db.Close()
}

//数据库查询COOKIE是否存在
func Mysqllcookiecx(r *http.Request) (bool,name){

	db:= MysqlMk()
	d:=name{}

	for _,k:=range r.Header["Cookie"]{

		query, _ := db.Query("SELECT * FROM `user`.`usercookie` WHERE `cookie` LIKE '" + k + "' LIMIT 0, 1")

		if query.Next() == true{
			err := query.Scan(&d.User, &d.Cookie)
			if err != nil {
				return false,d
			}
			return true,d
		}

	}

	return false,d

}

//数据库身份证查询
func Mysqlsfzcx(r *http.Request) (sfz){

	//连接数据库
	db:= MysqlMk()

	//数据结构
	d:=sfz{
	}
	if r.Form["sfz"] == nil {
		ff:=rand.Intn(100)
		r.Form["sfz"] = append(r.Form["sfz"],[]string{strconv.Itoa(ff),}...)
	}

	query, err := db.Query("SELECT * FROM `user`.`ND_sfz` WHERE `SFZ` LIKE '%"+r.Form["sfz"][0]+"%' LIMIT 0, 10")
	if err != nil {
		log.Println(err)
		return d
	}


	for i:=0;query.Next();i++ {
		var xm,sf [1]string
		query.Scan(&xm[0],&sf[0])
		//log.Println(xm,sf)
		d.Xm=append(d.Xm,xm[:]...)
		d.Sfz=append(d.Sfz,sf[:]...)
	}

	return d


}

