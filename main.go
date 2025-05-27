package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/settings"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go config.yaml")
		return
	}

	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Printf("Load config failed, err:%v\n", err)
		return
	}

	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("Init logger failed, err:%v\n", err)
		return
	}

	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()

	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("Init snowflake failed, err:%v\n", err)
		return
	}

	r := router.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("Run router failed, err:%v\n", err)
		return
	}
}
