package service

import (
	"log"
	"os"
	"time"

	"github.com/suryab-21/indico-test/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB function to init db
func InitDB() {
	if DB == nil {
		dsn := "host=pgsql_db user=postgres password=password dbname=pgsql_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
		config := &gorm.Config{
			Logger: logger.New(
				log.New(os.Stderr, "[GORM] ", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second,   // Slow SQL threshold
					LogLevel:                  logger.Silent, // Log level
					IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
					Colorful:                  true,          // Disable color
				},
			),
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			DisableForeignKeyConstraintWhenMigrating: true,
		}

		if db, err := gorm.Open(postgres.Open(dsn), config); err == nil {
			DB = db.Debug()
			AutoMigrate(DB)
		}
	}
}

func DBtest() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	AutoMigrate(db)

	DB = db
	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(ModelList...)
}

// ModelList list of model
var ModelList []interface{} = []interface{}{
	&model.User{},
	&model.Product{},
	&model.WarehouseLocation{},
	&model.Order{},
	&model.OrderItems{},
}
