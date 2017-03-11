package worker

import (
	"morego/golog"
	"morego/lib/antonholmquist/jason"
	"morego/lib/robfig/cron"
	//"fmt"
	"morego/area"
	"morego/global"
	z_type "morego/type"
	"os"
	"path/filepath"
	"strings"
	//"net"
)

// 获取服务器的根路径
func GetBase() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		golog.Error("GetBase Error ", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)

}

func GetConfig() *jason.Object {

	return global.ConfigJson

}

func GetEnableStatus() bool {
	if global.AppConfig.Enable <= 0 {
		return false
	} else {
		return true
	}

}

func Enable() bool {

	global.AppConfig.Enable = 1
	return true

}

func Disable() bool {

	global.AppConfig.Enable = 0
	return true

}

func AddCron(expression string, exefnc func()) bool {

	if cron, ok := global.Crons[expression]; ok {
		golog.Info("cron exist :", cron)
		return false
	}
	c := cron.New()
	c.AddFunc(expression, exefnc)
	c.Start()
	global.Crons[expression] = c
	return true

}

func RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true

}

func Get(key string) bool {

	return true

}

func Set(key string, value string) bool {

	return true

}

func GetSession(sid string) *z_type.Session {

	return nil

}

func Kick(sid string) bool {

	return true

}

func CreateChannel(id string, name string) bool {

	area.CreateChannel(id, name)
	return true
}

func RemoveChannel(id string) bool {

	return true
}

func GetChannels() bool {

	return true
}

func GetSidsByChannel(channel_id string) bool {

	return true

}

func ChannelAddSid(sid string, id string) bool {

	return true

}

func ChannelKickSid(sid string, id string) bool {

	return true

}

func Push(sid string, msg string) bool {

	return true

}

func PushBySids(sids []string, msg string) bool {

	return true

}

func PushAll(msg string) bool {

	return true

}

func Broadcast(msg string) bool {

	return true

}
