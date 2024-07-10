package schema

import (
	"log"

	"github.com/gocql/gocql"
)

func migrateUserTable(session *gocql.Session) {

	if err := session.Query(`
	 create table if not exists cme.tbl_User
	  ( id UUID primary key
		, name text
		, email text
		, phone text
		, password text
		, createdat timestamp ) ;
	
	`).Exec(); err != nil {
		log.Fatal("Unable to create table: ", err)
	}

	if err := session.Query(`
	 create index if not exists unique_column_index on cme.tbl_User (name) ;
	
	`).Exec(); err != nil {
		log.Fatal("Unable to create table: ", err)
	}

}

func migrateMessageTable(session *gocql.Session) {

	if err := session.Query(`
	 create table if not exists cme.tbl_Message
	  ( id UUID primary key
		, fromid UUID
		, toid UUID
		, messagetext text
		, createdat timestamp
		, readat timestamp ) ;
	
	`).Exec(); err != nil {
		log.Fatal("Unable to create table: ", err)
	}

}

func Migrate(session *gocql.Session) {

	migrateUserTable(session)
	migrateMessageTable(session)

}
