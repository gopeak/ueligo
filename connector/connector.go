package connector

import (
	"fmt"
	"net"
	"os"
	"time"
	"sync"
	"morego/global"
	"morego/golog"
)

var Glock *sync.Mutex
var ConnMlock *sync.RWMutex
var ChannelMlock *sync.RWMutex
var SessionMlock *sync.RWMutex
var UserChannelsMlock *sync.RWMutex




func checkError(err error) {
	if err != nil {
		golog.Error(os.Stderr, "Fatal error: %s", err.Error())
	}
}

func stat_tick() {

	timer := time.Tick(1000 * time.Millisecond)
	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		fmt.Println(time.Now().Unix(), " Connections: ", global.SumConnections, "  Qps: ", global.Qps)
	}
}

func user_tick(conn *net.TCPConn) {

	timer := time.Tick(5000 * time.Millisecond)
	for _ = range timer {
		ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }`, time.Now().Unix())
		go conn.Write([]byte(fmt.Sprintf("%s\r\n", ping)))
	}
}
