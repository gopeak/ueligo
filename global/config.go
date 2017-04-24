package global

import (
	"fmt"
	"math/rand"
	"morego/lib/BurntSushi/toml"
)

type configType struct {
	Name         string
	Enable       bool
	Status       string
	Version      string
	Loglevel     string
	RpcType      string
	PackType     string       `toml:"pack_type"`
	SingleMode   bool	  `toml:"single_mode"`
	Log          log          `toml:"log"`
	Admin        admin        `toml:"admin"`
	Connector    connector    `toml:"connector"`
	Object       object       `toml:"object"`
	WorkerServer workerServer `toml:"worker_server"`
	Hub          hub          `toml:"hub"`
	Area         area         `toml:"area"`
}

type log struct {
	LogLevel      string `toml:"log_level"`
	LogBehindType string `toml:"log_behind_type"`
	MongodbHost   string `toml:"mongodb_host"`
	MongodbPort   string `toml:"mongodb_port"`
}
type admin struct {
	HttpPort string `toml:"http_port"`
}

type connector struct {
	WebsocketPort     int `toml:"websocket_port"`
	SocketPort        int `toml:"socket_port"`
	MaxConections     int `toml:"max_conections"`
	MaxConntionsIp    int `toml:"max_conntions_ip"`
	MaxPacketRate     int `toml:"max_packet_rate"`
	MaxPacketRateUnit int `toml:"max_packet_rate_unit"`
	AuthCcmds	[]string `toml:"auth_cmds"`
}

type object struct {
	DataType      string `toml:"data_type"`
	RedisHost     string `toml:"redis_host"`
	RedisPort     string `toml:"redis_port"`
	RedisPassword string `toml:"redis_password"`
	MonogoHost    string `toml:"monogo_host"`
	MonogoPort    int    `toml:"3306"`
}

type workerServer struct {
	Servers [][]interface{} `toml:"servers"`
	ToHub []string  `toml:"to_hub"`
}

type hub struct {
	Hub_host string `toml:"hub_host"`
	Hub_port string `toml:"hub_port"`
}

type area struct {
	Init_area []string
}

var Config configType

func InitConfig() {

	if _, err := toml.DecodeFile("config.toml", &Config); err != nil {
		fmt.Println("toml.DecodeFile error:", err)
		return
	}

}

func GetRandWorkerAddr() string  {
	rand_index := rand.Intn(len(WorkerServers))
	return  WorkerServers[rand_index]
}

func InitWorkerAddr()   {

	for _,data := range Config.WorkerServer.Servers{
		worker_host, _ := data[0].(string)
		worker_port_str, _ := data[1].(string)
		WorkerServers = append( WorkerServers ,worker_host + ":" + worker_port_str )
	}
}

