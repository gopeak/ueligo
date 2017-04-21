package web

import (
	"fmt"
	"net/http"
	"os"
	"io"
 	"database/sql"
 	_"github.com/go-sql-driver/mysql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"time"
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
		db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/webim?timeout=90s&collation=utf8mb4_unicode_ci")
		if err != nil {
			resp:=fmt.Sprintf(format_str,500,"Open database error: "+err.Error() )
			w.Write([]byte(resp))
			return
		}
		defer db.Close()
		rows,err:= db.Query("select user,sign,sid,sign from user where user=? And pwd=? ",user,pwd)
		if err != nil {
			resp:=fmt.Sprintf(format_str,500,"Sql query error: "+err.Error() )
			w.Write([]byte(resp))
			return
		}
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for j := range values {
			scanArgs[j] = &values[j]
		}

		record := make(map[string]string)
		for rows.Next() {
			//将行数据保存到record字典
			err = rows.Scan(scanArgs...)
			for i, col := range values {
				if col != nil {
					record[columns[i]] = string(col.([]byte))
				}
			}
		}

		fmt.Println(record)
		json_encode ,err:=json.Marshal( record )
		resp := ""
		if _, ok := record[`user`]; ok {
			resp=fmt.Sprintf(`{ "code":%d ,"msg": "%s","data":%s} `,1,"验证成功",string(json_encode) )
		}else{
			resp=fmt.Sprintf(format_str,404,"用户名密码错误" )
		}
		w.Write([]byte(resp))
	}

}