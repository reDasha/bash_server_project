package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

var (
	Db  *sqlx.DB
	err error
)

type Command struct {
	ID     int    `db:"id" json:"id"`
	Script string `db:"script" json:"script"`
	Result string `db:"result" json:"result"`
}

func InitDB(cfg Config) {
	Db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")
}
