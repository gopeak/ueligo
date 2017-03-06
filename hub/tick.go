// main loop

package hub

import (
	"fmt"
	"morego/global"
	"morego/lib/syncmap"
	"time"
	//json2 "encoding/json"
	z_type "morego/type"

	"github.com/garyburd/redigo/redis"
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
	fmt.Println("reply", reply)
	if reply != nil {
		global.SyncUserSessions = reply.(*syncmap.SyncMap)
		var UserSessions = map[string]*z_type.Session{}
		for item := range global.SyncUserSessions.IterItems() {
			UserSessions[item.Key] = item.Value.(*z_type.Session)
			fmt.Println(UserSessions[item.Key].Sid)
		}
	}
}
