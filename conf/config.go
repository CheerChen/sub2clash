package conf

import "github.com/spf13/viper"

type Conf struct {
	Port   string
	Key    string
	Filter []string
	Subs   []string
}

var Cfg Conf

func Load(in string) (err error) {
	viper.SetConfigName(in)
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		//log.Fatalf("Error reading config file, %s", err)
		return
	}
	return viper.Unmarshal(&Cfg)
}
