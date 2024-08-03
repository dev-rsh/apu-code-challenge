package repository

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"service3/config"
	"sync"
	"time"
)

var (
	dbConn *gorm.DB
	dbOnce sync.Once
)

type baseRepository struct {
	dbConn *gorm.DB
}

func openDBConnection() *gorm.DB {
	dbConf := config.GetConfig().Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConf.Host,
		dbConf.Username,
		dbConf.Password,
		dbConf.DatabaseName,
		dbConf.Port,
		dbConf.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not establish connection to DB. Err is %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get db instance. Err is %s", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(15)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	return db
}

func GetDBConn() *gorm.DB {
	dbOnce.Do(func() {
		dbConn = openDBConnection()
	})

	return dbConn
}
