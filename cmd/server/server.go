package main

import (
	"fmt"
	"strconv"
	"time"

	rediscache "github.com/emersonary/go-authentication/cache/redis"
	"github.com/emersonary/go-authentication/config"
	httphandler "github.com/emersonary/go-authentication/controller/restapi"
	"github.com/emersonary/go-authentication/database/engine/cassandra"
	redisdb "github.com/emersonary/go-authentication/database/engine/redis"
	"github.com/emersonary/go-authentication/database/schema"
	"github.com/emersonary/go-authentication/webserver"
	"github.com/gocql/gocql"
)

func createDependencies(session *gocql.Session, config *config.Conf) {

	userHandler := httphandler.CreateUserDependency(session, config.TokenAuth)
	httphandler.CreateMessageDependency(session, config.TokenAuth, userHandler)

}

func waitFromStartUp(seconds int) {

	for i := 0; i <= seconds; i++ {

		time.Sleep(time.Second)

		if seconds-i != 0 {
			fmt.Println("Waiting for startup " + strconv.Itoa(seconds-i) + " second(s)")
		}
	}

}

func main() {

	config, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	waitFromStartUp(config.WaitFroStartUp)

	session, err := cassandra.CassandraDB(config)

	if err != nil {
		fmt.Println("Error connecting database:", err)
		return
	}

	redisSession, err := redisdb.RedisDB(config)

	if err != nil {
		fmt.Println("Error connecting database:", err)
		return
	}

	rediscache.RedisCtrl = rediscache.NewRedisControl(redisSession)

	schema.Migrate(session)

	webserver.InitMiddlewares(config)

	createDependencies(session, config)

	webserver.StartWebServer(config)

}
