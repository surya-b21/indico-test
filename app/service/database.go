package service

import (
	"log"
	"os"
	"time"

	"github.com/suryab-21/indico-test/app/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB function to init db
func InitDB() {
	if DB == nil {
		dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
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

		if db, err := gorm.Open(mysql.Open(dsn), config); err == nil {
			DB = db.Debug()
			AutoMigrate(DB)
		}
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(ModelList...)
}

// ModelList list of model
var ModelList []interface{} = []interface{}{
	&model.User{},
}
