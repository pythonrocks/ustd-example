package main

import (
	"log"

	"github.com/spf13/viper"

	server "github.com/pythonrocks/ustd-example/internal/api"
	"github.com/pythonrocks/ustd-example/internal/conf"
	"github.com/pythonrocks/ustd-example/internal/service"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var conf conf.Conf

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Println("Initializing environment...")
	env := &service.Env{
		Conf: &conf,
	}

	log.Println("Starting RPC server on port", conf.RPCPort)
	if err := server.StartRPCServer(env); err != nil {
		log.Fatal(err)
	}

}
