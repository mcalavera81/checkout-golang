package config

import "log"
import (
	"github.com/kelseyhightower/envconfig"
	"strings"
	"strconv"
	"github.com/sirupsen/logrus"
)
var config ClientConfig


/*type Specification struct {
	Debug       bool
	Port        int
	User        string
	Users       []string
	Rate        float32
	Timeout     time.Duration
	ColorCodes  map[string]int
}
*/

func GetConfig() ClientConfig {
	return config
}

type ClientConfig struct {
	Host		    string `default:"localhost"`
	Port		    int 	`default:"9000" envconfig:"PORT"`
	BaseURI		    string `ignored:"true"`
	CheckoutPath	string `default:"/api/checkout/"`
	LogLevel logrus.Level
}



func init(){

	err := envconfig.Process("client", &config)
	config.BaseURI=strings.Join([]string{"http://",config.Host,":",strconv.Itoa(config.Port),config.CheckoutPath},"")
	if err != nil {
		log.Fatal(err.Error())
	}


}