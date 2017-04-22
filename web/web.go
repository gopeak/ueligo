package web

import (
	"fmt"
	"net/http"
	"os"
	"io"
 	_"github.com/go-sql-driver/mysql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"time"
	"strconv"
	"database/sql"
	"log"
)



func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data": {"src":"%s","name":"%s" }} `

	if r.Method == "GET" {
		resp:=fmt.Sprintf(format_str,401,"GET no support!","","")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			//fmt.Println(err)
			resp:=fmt.Sprintf(format_str,400,err.Error(),"","")
			w.Write([]byte(resp))
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "%v", handler.Header)
		wd, _ := os.Getwd()
		upload_dir := fmt.Sprintf("%s/web/wwwroot/data/images/", wd)
		f, err := os.OpenFile(upload_dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		code := 0
		err_str := ""
		src:="http://"+r.Host+"/data/images/"+handler.Filename
		if err != nil {
			fmt.Println(err)
			code = 500
			err_str = err.Error()
		}else{
			defer f.Close()
			io.Copy(f, file)
		}
		resp:=fmt.Sprintf(format_str,code,err_str,src,handler.Filename,)
		w.Write([]byte(resp))
	}

}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data": {"src":"%s","name":"%s" }} `

	if r.Method == "GET" {
		resp:=fmt.Sprintf(format_str,401,"GET no support!","","")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			//fmt.Println(err)
			resp:=fmt.Sprintf(format_str,400,err.Error(),"","")
			w.Write([]byte(resp))
			return
		}
		defer file.Close()

		//fmt.Fprintf(w, "%v", handler.Header)
		wd, _ := os.Getwd()
		upload_dir := fmt.Sprintf("%s/web/wwwroot/data/files/", wd)
		f, err := os.OpenFile(upload_dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		code := 0
		err_str := ""
		src:="http://"+r.Host+"/data/files/"+handler.Filename
		if err != nil {
			fmt.Println(err)
			code = 500
			err_str = err.Error()
		}else{
			defer f.Close()
			io.Copy(f, file)
		}
		resp:=fmt.Sprintf(format_str,code,err_str,src,handler.Filename,)
		w.Write([]byte(resp))
	}

}


func RegHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data": { }} `

	if r.Method == "GET" {
		resp:=fmt.Sprintf(format_str,401,"GET no support!" )
		w.Write([]byte(resp))
		return

	} else {

		r.ParseForm( )

		user  := r.PostForm.Get(`user`)
		pwd  := r.PostForm.Get(`pwd`)
		age  := r.PostForm.Get(`age`)
		nick  := r.PostForm.Get(`nick`)
		sign  := r.PostForm.Get(`sign`)
		reg_time  := time.Now().Unix()

		db:=new(Mysql)
		db.Connect()

		row := db.GetRow( `select user from user where user=? `,user)
		if _, ok := row[`user`]; ok {
			resp:=fmt.Sprintf(format_str,0 ,"用户名已经存在!")
			w.Write([]byte(resp))
			return
		}

		insertid,err:=db.Insert( `INSERT user (user,pwd,nick,sign,age,reg_time) values (?,?,?,?,?,?)` ,user,pwd,nick,sign,age,reg_time)
		if( err!=nil ){
			resp:=fmt.Sprintf(format_str,500,"db.Insert err:",err.Error() )
			w.Write([]byte(resp))
			return
		}
		fmt.Println( "insertid:", insertid )
		resp:=fmt.Sprintf(format_str,1,"注册成功" )
		w.Write([]byte(resp))
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data": {  }} `

	if r.Method == "GET" {
		resp:=fmt.Sprintf(format_str,401,"GET no support!" )
		w.Write([]byte(resp))
		return

	} else {
		r.ParseForm( )
		user  := r.PostForm.Get(`user`)
		pwd  := r.PostForm.Get(`pwd`)
		fmt.Println( user ,pwd ,mysql.MySQLDriver{})
		db:=new(Mysql)
		db.Connect()

		record:= getUserRow(db.Db,  "select id,user,sign,sid from user where user=? and pwd=?",user,pwd)
		fmt.Println(record)
		json_encode ,_:=json.Marshal( record )
		resp := ""
		if _, ok := record[`user`]; ok {
			resp=fmt.Sprintf(`{ "code":%d ,"msg": "%s","data":%s} `,1,"验证成功",string(json_encode) )
		}else{
			resp=fmt.Sprintf(format_str,404,"用户名密码错误" )
		}
		w.Write([]byte(resp))
	}

}
func   getUserRow(db *sql.DB, sql_str string, args ...interface{})  map[string]string {

	var  id,user,sign,sid,avatar string
	record := make(map[string]string)
	err := db.QueryRow(sql_str, args...).Scan(&id,&user,&sign,&sid,&avatar)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("id is %s\n", id)
	}
	record["id"] = id
	record["user"] = user
	record["sign"] = sign
	record["sid"] = sid
	record["avatar"] = avatar

	return record
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data":%s} `

	if r.Method == "GET" {

		r.ParseForm( )
		id_str  := r.Form.Get(`id`)
		id ,_ := strconv.Atoi( id_str )
		sid  := r.Form.Get(`sid`)
		fmt.Println( id ,sid ,mysql.MySQLDriver{})
		db:=new(Mysql)
		_,err:=db.Connect()
		if( err!=nil ){
			resp:=fmt.Sprintf(format_str,500,"db.Insert err:",err.Error() ,"{}")
			w.Write([]byte(resp))
			return
		}

		// 获取当前用户信息
		sql := "select id,nick as username,status ,sign,avatar  from user where sid=? "
		my_record:= db.GetRow( sql, sid )
		resp := ""
		_, ok := my_record[`user`];
		if !ok {
			resp=fmt.Sprintf(format_str,404,"验证错误" )
			w.Write([]byte(resp))
			return
		}
		uid ,_:= strconv.Atoi( my_record["id"] )

		// 获取所属的联系人列表（未分组）
		sql = "SELECT  u.id,u.nick as username,u.avatar,u.sign,c.group_id FROM `contacts` c LEFT JOIN `user` u on u.id =c.uid WHERE  c.master_uid=?"
		contact_records := db.GetRows(sql,uid)
		fmt.Println( contact_records )

		// 获取分组
		sql = "SELECT  id,title as groupname  FROM `contact_group` WHERE uid=? "
		my_group_records := db.GetRows( sql ,uid)

		for _key,group  :=  range my_group_records{
			group_id := group[`id`]
			tmp_list := make([]map[string]string,0)

			for _k,c  :=  range contact_records{
				if c[`group_id`]==group_id{
					tmp_list = append( tmp_list , c )
					contact_records=append(contact_records[:_k],contact_records[_k+1:]...)
				}
			}
			tmp_list_str,_ := json.Marshal( tmp_list )
			group["list"] = string(tmp_list_str)
			group["online"] = "1"
			my_group_records[_key] = group

		}
		fmt.Println( my_group_records )

		// 获取群组
		sql = "SELECT id,channel_id,pic as avatar FROM `global_group` WHERE  id in( SELECT `group_id` FROM `user_join_group` WHERE `uid`=? )"
		join_group_records := db.GetRows( sql,uid )
		fmt.Println( join_group_records )

		data_format_str := `{ "mine":%s ,"friend": "%s","group":%s} `
		my_record_encode ,err:=json.Marshal( my_record )
		my_group_records_encode ,err:=json.Marshal( my_group_records )

		join_group_records_encode ,err:=json.Marshal( join_group_records )
		data_resp := fmt.Sprintf(data_format_str,string(my_record_encode),string(my_group_records_encode) ,string(join_group_records_encode))
		fmt.Println( data_resp )
		resp=fmt.Sprintf(format_str,1,data_resp )
		w.Write([]byte(resp))
	}

}