package cassandra

import (
	"fmt"
	"strconv"

	"github.com/emersonary/go-authentication/config"
	"github.com/gocql/gocql"
)

func CassandraDB(config *config.Conf) (*gocql.Session, error) {

	cassandraDB := gocql.NewCluster(config.DBHost)
	cassandraDB.Keyspace = config.DBName
	cassandraDB.Port = config.DBPort
	cassandraDB.Authenticator = gocql.PasswordAuthenticator{Username: config.DBUser, Password: config.DBPassword}
	cassandraDB.Consistency = gocql.Quorum
	session, err := cassandraDB.CreateSession()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to Cassandra on " + config.DBHost + ":" + strconv.Itoa(config.DBPort) + ", keyspace " + cassandraDB.Keyspace)
	return session, nil

}
