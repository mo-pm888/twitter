package pg

import (
	"database/sql"
	"fmt"

	"Twitter_like_application/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type ServicePostgresql struct {
	DB *sql.DB
}

func ConnectPostgresql(c config.Config) error {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("error connect BD")
		return err
	}

	DB = db
	fmt.Println("**** PG ran.... ****")

	return nil
}
