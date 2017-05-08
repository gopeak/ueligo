// main loop

package hub

import (
	"fmt"
	"morego/global"
	"morego/lib/syncmap"
	"time"
	//json2 "encoding/json"
	"morego/golog"
	z_type "morego/type"
	"net"

	"github.com/garyburd/redigo/redis"
	//"morego/protocol"
)

func tick() {
	timer := time.Tick(100 * time.Millisecond)
	for now := range timer {
		// entity updates (you could use now for physic engine calculs)
		// this is called every 100 millisecondes
		// playerFactory.Update()
		fmt.Println("now", now)
	}
}

func TickSyncSession() {

	redisc, err := redis.Dial("tcp", global.Config.Object.RedisHost+`:`+string(global.Config.Object.RedisPort))
	//defer redisc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	timer := time.Tick(1 * time.Second)
	var LastSessions *syncmap.SyncMap

	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		/*var UserSessions = map[string]*z_type.Session{}
		for item := range global.SyncUserSessions.IterItems() {
			UserSessions[item.Key] = item.Value.(*z_type.Session)
		}
		js1, _ := json2.Marshal(UserSessions)
		*/
		if LastSessions != global.SyncUserSessions {
			redisc.Do("Set", "morego/user_session", global.SyncUserSessions)
			redisc.Flush()
			LastSessions = global.SyncUserSessions
		}

	}
}

func LoadSessionFromRedis() {

	redisc, err := redis.Dial("tcp", global.Config.Object.RedisHost+`:`+string(global.Config.Object.RedisPort))
	//defer redisc.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	reply, err_get := redisc.Do("Get", "morego/user_session")
	if err_get != nil {
		fmt.Println(err_get)
		return
	}
	fmt.Println("GET morego/user_session ", reply)
	if reply != nil {
		global.SyncUserSessions = reply.(*syncmap.SyncMap)
		var UserSessions = map[string]*z_type.Session{}
		for item := range global.SyncUserSessions.IterItems() {
			UserSessions[item.Key] = item.Value.(*z_type.Session)
			fmt.Println(UserSessions[item.Key].Sid)
		}
	}
}

func TickWorkerServer() {
	// 先暂停10秒
	time.Sleep(5 * time.Second)
	timer := time.Tick(10 * time.Second)
	for now := range timer {
		//fmt.Println("now", now)
		ch_success := make(chan string, 0)
		for _, data := range global.Config.WorkerServer.Servers {
			go func(data []interface{}) {
				worker_host, _ := data[0].(string)
				worker_port_str, _ := data[1].(string)
				ip_port := worker_host + ":" + worker_port_str

				//fmt.Println("tcpAddr: ",index," ", ip_port)
				conn, err_req := net.DialTimeout("tcp", ip_port, 5*time.Second)
				if err_req != nil {
					golog.Error("检测到 workerserver:", ip_port, " 连接异常!", now)
					for i, addr := range global.WorkerServers {
						if addr == ip_port {
							global.WorkerServers = append(global.WorkerServers[:i], global.WorkerServers[i+1:]...)
						}
					}
					ch_success <- ip_port + err_req.Error()
				} else {
					exist := false
					for _, addr := range global.WorkerServers {
						if addr == ip_port {
							exist = true
							break
						}
					}
					if !exist {
						global.WorkerServers = append(global.WorkerServers, ip_port)
					}
					ch_success <- ip_port + "ok"
				}
				//fmt.Println("result: ", ip_port, " ok")
				//req_str:= fmt.Sprintf("%d||%s||%s||%d||%s\n", protocol.TypePing, "Ping", "", 0, "")
				//conn.Write([]byte(req_str))
				conn.Close()
			}(data)
		}
		sum := 0
		for i := 0; i < len(global.Config.WorkerServer.Servers)+1; i++ {
			select {
			case <-ch_success:
				//fmt.Println("recv_result:", r)
				sum++
				if sum == len(global.Config.WorkerServer.Servers) {
					break
				}

			default:
				//fmt.Printf(".")
				time.Sleep(10 * time.Millisecond)
			}
		}
		//fmt.Println("sum:", sum)

	}
}
