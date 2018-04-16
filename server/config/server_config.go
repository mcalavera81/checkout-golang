package serverconfig

import (
	"github.com/kelseyhightower/envconfig"
	"time"
	log "github.com/sirupsen/logrus"
)
var config = new(ServerConfig)



func GetConfig() *ServerConfig {
	return config
}

type ServerConfig struct {
	Port	int `default:"9000" envconfig:"PORT"`
	Timeout time.Duration `default:10`
	LogLevel LogLevel `default:"info" envconfig:"LOG_LEVEL"`
}

type LogLevel log.Level

func (ll *LogLevel) Decode(value string) error {
	level, err := log.ParseLevel(value)
	if  err != nil{
		return err
	}

	*ll = LogLevel(level)
	return nil

}

func init(){

	err := envconfig.Process("server", config)
	if err != nil {
		log.Fatal(err.Error())
	}


}