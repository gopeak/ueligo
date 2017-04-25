package web

import (
	_"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"morego/area"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": {"src":"%s","name":"%s" }} `

	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!", "", "")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			//fmt.Println(err)
			resp := fmt.Sprintf(format_str, 400, err.Error(), "", "")
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
		src := "http://" + r.Host + "/data/images/" + handler.Filename
		if err != nil {
			fmt.Println(err)
			code = 500
			err_str = err.Error()
		} else {
			defer f.Close()
			io.Copy(f, file)
		}
		resp := fmt.Sprintf(format_str, code, err_str, src, handler.Filename)
		w.Write([]byte(resp))
	}

}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": {"src":"%s","name":"%s" }} `

	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!", "", "")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			//fmt.Println(err)
			resp := fmt.Sprintf(format_str, 400, err.Error(), "", "")
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
		src := "http://" + r.Host + "/data/files/" + handler.Filename
		if err != nil {
			fmt.Println(err)
			code = 500
			err_str = err.Error()
		} else {
			defer f.Close()
			io.Copy(f, file)
		}
		resp := fmt.Sprintf(format_str, code, err_str, src, handler.Filename)
		w.Write([]byte(resp))
	}

}

func RegHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	format_str := `{ "code":%d ,"msg": "%s","data": { }} `

	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!")
		w.Write([]byte(resp))
		return

	} else {

		r.ParseForm()

		user := r.PostForm.Get(`user`)
		pwd := r.PostForm.Get(`pwd`)
		age := r.PostForm.Get(`age`)
		nick := r.PostForm.Get(`nick`)
		sign := r.PostForm.Get(`sign`)
		avatar := r.PostForm.Get(`avatar`)
		reg_time := time.Now().Unix()
		sid := area.CreateSid()

		db := new(Mysql)
		db.Connect()

		row := db.GetRow(`select user from user where user=? `, user)

		if _, ok := row[`user`]; ok {
			resp := fmt.Sprintf(format_str, 0, "用户名已经存在!")
			w.Write([]byte(resp))
			return
		}

		insert_id, err := db.Insert(`INSERT user (user,pwd,nick,sign,age,sid,avatar,reg_time)
						    values (?,?,?,?,?,?,?,?)`,
			user, pwd, nick, sign, age,sid,avatar, reg_time)
		if err != nil {
			resp := fmt.Sprintf(format_str, 500, "db.Insert err:", err.Error())
			w.Write([]byte(resp))
			return
		}
		fmt.Println("insertid:", insert_id)
		resp := fmt.Sprintf(format_str, 1, "注册成功")
		w.Write([]byte(resp))
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s", "data": {}} `
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if r.Method == "GET" {
		resp := fmt.Sprintf(format_str, 401, "GET no support!")
		w.Write([]byte(resp))
		return

	} else {
		r.ParseForm()
		user := r.PostForm.Get(`user`)
		pwd := r.PostForm.Get(`pwd`)
		fmt.Println(user, pwd, mysql.MySQLDriver{})
		db := new(Mysql)
		db.Connect()

		resp := ""
		sql_str := `select id,user,sign,sid ,avatar from user  where user=? and pwd=?`
		var id, sign, sid, avatar string
		record := make(map[string]string)
		scan_err := db.Db.QueryRow(sql_str, user, pwd).Scan(&id, &user, &sign, &sid, &avatar)
		if scan_err != nil {
			resp = fmt.Sprintf(format_str, 500, "用户名密码错误"+scan_err.Error())
			w.Write([]byte(resp))
			return
		}
		record["id"] = id
		record["user"] = user
		record["sign"] = sign
		record["sid"] = sid
		record["avatar"] = avatar
		token := area.CreateSid()
		affect_num,_:=db.Update( `Update user set token=? Where id=?`,token,id)
		if affect_num>0 {
			record["token"] = token
		}

		fmt.Println(record)
		json_encode, _ := json.Marshal(record)

		uid, _ := strconv.Atoi(record["id"])
		friends := getMyContacts(db.Db, uid)
		friends_encode, _ := json.Marshal(friends)

		if id != "" {
			resp = fmt.Sprintf(`{ "code":%d ,"msg": "%s","data":%s,"contacts":%s} `, 1, "验证成功", string(json_encode) ,string(friends_encode))
		} else {
			resp = fmt.Sprintf(format_str, 404, "用户名密码错误")
		}
		w.Write([]byte(resp))
	}

}


func GetListHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)
	_list := new(ListType)
	root.Data = &_list

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		r.ParseForm()
		id_str := r.Form.Get(`id`)
		id, _ := strconv.Atoi(id_str)
		sid := r.Form.Get(`sid`)
		fmt.Println(id, sid, mysql.MySQLDriver{})
		db := new(Mysql)
		_, err := db.Connect()
		if err != nil {
			root.Code = 500
			root.Msg = "数据库连接失败:" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}
		uid, _ := strconv.Atoi(my_record["id"])
		_list.Mine = my_record
		_list.Friend = getFriends(db.Db, uid)
		_list.Group = getMyGroups( db.Db,uid)

		root.Code = 0
		root.Msg = ""
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}
}


func GetMemberHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	root := new(Root)
	member := new(MemberType)
	root.Data = &member

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		r.ParseForm()
		id_str := r.Form.Get(`id`)
		id, _ := strconv.Atoi(id_str)
		sid := r.Form.Get(`sid`)
		fmt.Println(id, sid, mysql.MySQLDriver{})
		db := new(Mysql)
		_, err := db.Connect()
		if err != nil {
			root.Code = 500
			root.Msg = "数据库连接失败:" + err.Error()
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		// 获取当前用户信息
		my_record := GetUserRow(db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			root.Code = 400
			root.Msg = "用户验证失败"
			json_encode ,_:=json.Marshal( root )
			w.Write( json_encode )
			return
		}

		member.Owner = my_record
		member.List = getMembers(db.Db, id)
		member.Members = len( member.List  )

		root.Code = 0
		root.Msg = ""
		json_encode ,_:=json.Marshal( root )
		w.Write( json_encode )
	}

}

