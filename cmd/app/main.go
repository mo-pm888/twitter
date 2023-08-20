package main

import (
	"fmt"
	"log"

	"Twitter_like_application/config"
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/server"
	"Twitter_like_application/migrations"
)

type ServiceMongoDb struct {
	DB interface{}
}

func main() {
	c, err := config.New()
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
	err = pg.ConnectPostgresql(*c)
	if err != nil {
		fmt.Println(err)
	}
	err = migrations.Run(pg.DB)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("**** running migrations ****", err)
	}
	err = server.Server(*c)
	if err != nil {
		fmt.Println(err)
	}

}
