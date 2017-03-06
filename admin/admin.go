package admin

import (
	"encoding/json"
	"fmt"
	"morego/global"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	cpu "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type MongoLog struct {
	Id_     bson.ObjectId `bson:"_id"`
	Name    string
	Level   string
	File    string
	Line    int
	Message string
	Time    int
}

type Logs struct {
	All []MongoLog
}

var (
	mgoSession *mgo.Session
	dataBase   = "gomore"
)

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(global.Config.Log.MongodbHost)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(dataBase).C(collection)
	return s(c)
}

func getCollection(collection string) *mgo.Collection {
	session := getSession()
	c := session.DB(dataBase).C(collection)
	return c
}

/**
 * 建立web server
 */
func HttpServer() {

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/admin/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))
	http.HandleFunc("/stats", statsTask)
	http.HandleFunc("/lastlogs", lastLogsTask)
	http.HandleFunc("/searchLogs", searchLogsTask)

	go func() {
		http.ListenAndServe(":"+global.Config.Admin.HttpPort, nil)
	}()
}

func statsTask(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("statsTask is running...")
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct
	//fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)
	cpuf, _ := cpu.Percent(1, true)
	//fmt.Println("Percent:", cpuf)

	now := time.Now().Format("15:04:05")
	str := fmt.Sprintf(`{"time":"%s","conns":%d,"qps":%d,"cpu_per":%v,"mem_total":"%v","mem_free":"%v" , "mem_use_per":"%f"}`,
		now, global.SumConnections, global.Qps, cpuf, v.Total, v.Free, v.UsedPercent)
	w.Write([]byte(str))

}

func lastLogsTask(w http.ResponseWriter, req *http.Request) {
	fmt.Println("lastLogsTask is running...")

	session := getSession()
	defer session.Close()

	collection := getCollection("logs")

	ms := []MongoLog{}
	err := collection.Find(bson.M{}).All(&ms)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("All:", ms)
	if b, err := json.Marshal(ms); err == nil {
		w.Write(b)
	} else {
		w.Write([]byte(`[]`))
	}

}

func searchLogsTask(w http.ResponseWriter, req *http.Request) {

	fmt.Println("searchLogsTask is running...")

	session := getSession()
	defer session.Close()

	results, _ := SearchLog("logs", bson.M{}, `Time`, bson.M{}, 0, 100)

	if b, err := json.Marshal(results); err == nil {
		w.Write(b)
	} else {
		w.Write([]byte(`[]`))
	}

}

/**
 * 执行查询，此方法可拆分做为公共方法
 * [Search description]
 * @param {[type]} collectionName string [description]
 * @param {[type]} query          bson.M [description]
 * @param {[type]} sort           bson.M [description]
 * @param {[type]} fields         bson.M [description]
 * @param {[type]} skip           int    [description]
 * @param {[type]} limit          int)   (results      []interface{}, err error [description]
 */
func SearchLog(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
	exop := func(c *mgo.Collection) error {
		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
	}
	err = witchCollection(collectionName, exop)
	return
}
