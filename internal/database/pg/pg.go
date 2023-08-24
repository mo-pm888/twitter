package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"Twitter_like_application/config"

	_ "github.com/lib/pq"
)

const (
	pgGoodMsg = "postgres starts"
	pgBadMsg  = "error connect BD"
)

func ConnectPostgresql(c config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.New(pgBadMsg)
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println(pgGoodMsg)
	return db, nil
}
