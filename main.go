package main

import (
	"database/sql"
	"log"

	"github.com/Franklynoble/bankapp/api"
	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/Franklynoble/bankapp/db/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".") // this Load would come from current file

	if err != nil {
		log.Fatal("cannot loadss config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("can not connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

}
