package web

import (
	"fmt"
	"net/http"
	"os"
	"io"
 	"database/sql"
 	"github.com/go-sql-driver/mysql"
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

		resp:=fmt.Sprintf(format_str,2,"" )
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
		row:= db.QueryRow("select user,pwd from user where user=? And pwd=? ",user,pwd)

		var row_user string
		var row_pwd string

		err = row.Scan(&row_user, &row_pwd)
		if err != nil {
			resp:=fmt.Sprintf(format_str,500,err.Error() )
			w.Write([]byte(resp))
			return
		}
		log.Println(row_user, row_pwd)

		resp:=fmt.Sprintf(format_str,1,err.Error() )
		w.Write([]byte(resp))
	}

}