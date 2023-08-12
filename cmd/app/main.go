package main

import (
	"Twitter_like_application/internal/admin"
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/server"
	"Twitter_like_application/migrations"
	"fmt"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	err := admin.LoadEnvFile()
	if err != nil {
		fmt.Println(err)
	}
	err = pg.ConnectPostgresql()
	if err != nil {
		fmt.Println(err)
	}
	err = migrations.Run(pg.DB)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("**** running migrations ****", err)
	}
	err = server.Server()
	if err != nil {
		fmt.Println(err)
	}

}
