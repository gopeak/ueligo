package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

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
		reg_time := time.Now().Unix()

		db := new(Mysql)
		db.Connect()

		row := db.GetRow(`select user from user where user=? `, user)

		if _, ok := row[`user`]; ok {
			resp := fmt.Sprintf(format_str, 0, "用户名已经存在!")
			w.Write([]byte(resp))
			return
		}

		insertid, err := db.Insert(`INSERT user (user,pwd,nick,sign,age,reg_time) values (?,?,?,?,?,?)`, user, pwd, nick, sign, age, reg_time)
		if err != nil {
			resp := fmt.Sprintf(format_str, 500, "db.Insert err:", err.Error())
			w.Write([]byte(resp))
			return
		}
		fmt.Println("insertid:", insertid)
		resp := fmt.Sprintf(format_str, 1, "注册成功")
		w.Write([]byte(resp))
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data": {  }} `

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
		scan_err := db.Db.QueryRow(sql_str, user, pwd ).Scan(&id, &user, &sign, &sid, &avatar)
		if( scan_err!=nil ){
			resp = fmt.Sprintf(format_str, 500, "用户名密码错误"+scan_err.Error())
			w.Write([]byte(resp))
			return
		}
		record["id"] = id
		record["user"] = user
		record["sign"] = sign
		record["sid"] = sid
		record["avatar"] = avatar
		fmt.Println("id:",id)
		fmt.Println(record)
		json_encode, _ := json.Marshal(record)

		if id!="" {
			resp = fmt.Sprintf(`{ "code":%d ,"msg": "%s","data":%s} `, 1, "验证成功", string(json_encode))
		} else {
			resp = fmt.Sprintf(format_str, 404, "用户名密码错误")
		}
		w.Write([]byte(resp))
	}

}
func getUserRow(db *sql.DB,sid string ) map[string]string {

	sql_str := `select id,nick,status ,sign, avatar  from user where sid=?`
	var id, nick,status, sign, avatar string
	record := make(map[string]string)
	err := db.QueryRow(sql_str,sid ).Scan(&id, &nick, &status, &sign,  &avatar)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("id is %s\n", id)
	}
	record["id"] = id
	record["username"] = nick
	record["sign"] = sign
	record["status"] = status
	record["sid"] = sid
	record["avatar"] = avatar

	return record
}

func GetListHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method) //获取请求的方法

	format_str := `{ "code":%d ,"msg": "%s","data":%s } `
	record := make(map[string]string)
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
			resp := fmt.Sprintf(format_str, 500, "db.Insert err:", err.Error(), "{}")
			w.Write([]byte(resp))
			return
		}

		resp := ""
		// 获取当前用户信息
		my_record := getUserRow( db.Db, sid)
		_, ok := my_record[`id`]
		if !ok {
			resp = fmt.Sprintf(format_str, 404, "验证错误")
			w.Write([]byte(resp))
			return
		}
		uid, _ := strconv.Atoi(my_record["id"])

		// 获取所属的联系人列表（未分组）
		sql := "SELECT  u.id,u.nick as nick,u.avatar,u.sign,c.group_id FROM `contacts` c LEFT JOIN `user` u on u.id =c.uid WHERE  c.master_uid=?"

		contact_records := make([]map[string]string,0)
		rows, err := db.Db.Query(sql,uid)
		if err != nil {
			resp = fmt.Sprintf(format_str, 501, "服务器错误@"+err.Error())
			w.Write([]byte(resp))
			return
		}
		for rows.Next() {
			//将行数据保存到record字典
			var id, nick, avatar, sign,group_id string
			record := make(map[string]string)
			err = rows.Scan( &id, &nick, &avatar, &sign, &group_id )
			if( err!=nil ){
				resp = fmt.Sprintf(format_str, 502 , "服务器错误@"+err.Error())
				w.Write([]byte(resp))
				return
			}
			record["id"] = id
			record["username"] = nick
			record["avatar"] = avatar
			record["sign"] = sign
			record["group_id"] = group_id
			fmt.Println(record)
			contact_records = append( contact_records,record )

		}
		fmt.Println(contact_records)

		// 获取分组
		sql = "SELECT  id,title as groupname  FROM `contact_group` WHERE uid=? "
		my_group_records := make([]map[string]string,0)
		rows, err = db.Db.Query(sql,uid)
		if err != nil {
			resp = fmt.Sprintf(format_str, 504, "服务器错误@"+err.Error())
			w.Write([]byte(resp))
			return
		}
		for rows.Next() {
			//将行数据保存到record字典
			var gid, groupname  string
			err = rows.Scan( &gid, &groupname )
			if( err!=nil ){
				resp = fmt.Sprintf(format_str, 505, "服务器错误@"+err.Error())
				w.Write([]byte(resp))
				return
			}
			record["id"] = gid
			record["groupname"] = groupname
			fmt.Println(record)
			my_group_records = append( my_group_records,record )
		}

		for _, group := range my_group_records {
			group_id := group[`id`]
			tmp_list := make([]map[string]string, 0)

			for _k, c := range contact_records {
				if c[`group_id`] == group_id {
					tmp_list = append(tmp_list, c)
					contact_records = append(contact_records[:_k], contact_records[_k+1:]...)
				}
			}
			tmp_list_str, _ := json.Marshal(tmp_list)
			group["list"] = string(tmp_list_str)
			group["online"] = "1"
			//my_group_records[_key] = fmt.Sprintf(`{ "groupname": "%s","id": "%s","list": [],"online": "1"}`,group["groupname"],group["id"])

		}
		fmt.Println(my_group_records)

		// 获取群组
		sql = "SELECT id,channel_id,pic as avatar FROM `global_group` WHERE  id in( SELECT `group_id` FROM `user_join_group` WHERE `uid`=? )"
		join_group_records := make([]map[string]string,0)
		rows, err = db.Db.Query(sql,uid)
		if err != nil {
			resp = fmt.Sprintf(format_str, 504, "服务器错误@"+err.Error())
			w.Write([]byte(resp))
			return
		}
		for rows.Next() {
			//将行数据保存到record字典
			var cid, channel_id ,avatar string
			err = rows.Scan( &cid, &channel_id,&avatar )
			if( err!=nil ){
				resp = fmt.Sprintf(format_str, 505, "服务器错误@"+err.Error())
				w.Write([]byte(resp))
				return
			}
			record["id"] = cid
			record["channel_id"] = channel_id
			record["avatar"] = avatar
			fmt.Println(record)
			join_group_records = append( join_group_records,record )
		}
		fmt.Println(join_group_records)

		data_format_str := `{ "mine":%s ,"friend": %s,"group":%s} `
		my_record_encode, err := json.Marshal(my_record)
		my_group_records_encode, err := json.Marshal(my_group_records)

		join_group_records_encode, err := json.Marshal(join_group_records)
		data_resp := fmt.Sprintf(data_format_str, string(my_record_encode), string(my_group_records_encode), string(join_group_records_encode))
		fmt.Println(data_resp)
		resp = fmt.Sprintf(format_str, 0,"ok", data_resp)

		w.Write([]byte(resp))
	}

}
