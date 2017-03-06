/**
 *  结构体定义
 *
 */

package global

import (

	//"code.google.com/p/go.net/websocket"
	//"net"
	"morego/lib/syncmap"
)

type Session struct {
	IP string
	// user json data
	User string
	// session related
	LoggedIn bool // flag for weather the user is logged in
	KickOut  bool // flag for player is kicked out
	// session flag
	Sid string
	// time related variables
	ConnectTime int64 // tcp connection establish time, in millsecond(ms)
	PacketTime  int64 // last packet time
}

type Appconfig struct {
	Enable int64 //  Listen to clients
	Status string
}

type UserChannelArrs struct {
	Names []string //  List of ready workers
}

type ChannelRpcType struct {
	Id         string
	Name       string
	Host       string
	Port       string
	Conns      *syncmap.SyncMap
	WsConns    *syncmap.SyncMap
	CreateTime int64
}
