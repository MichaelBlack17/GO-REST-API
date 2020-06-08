package main

import (
	"GO-REST-API/internal/app/apiserver"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init(){
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml","path to config file")
}

func main()  {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil{
		logrus.Fatal(err)
	}


	if err := apiserver.Start(config); err != nil {
	logrus.Fatal(err)
	}
}
