package database

import "gorm.io/gorm"

type DatabasesPostgres interface {
	Connect() *gorm.DB
}
