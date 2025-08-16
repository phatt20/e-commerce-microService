package database

import (
	"fmt"
	"log"
	"microService/config"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	*gorm.DB
}

var (
	postgrestDatabaseInstance *postgresDatabase
	once                      sync.Once
)

func NewPostgresDatabase(conf *config.PostgresConfig) DatabasesPostgres {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
			conf.Host, conf.Port, conf.User, conf.Password, conf.DBName, conf.SSLMode, conf.Schema,
		)
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		log.Printf("Connected to database %s", conf.DBName)

		postgrestDatabaseInstance = &postgresDatabase{conn}
	})
	return postgrestDatabaseInstance
}

func (db *postgresDatabase) Connect() *gorm.DB {
	return db.DB
}
