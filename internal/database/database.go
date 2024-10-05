package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service interface {

	// Close terminates the database connection.
	// It returns an error if connection cannot be closed.
	Close() error

	// Return active GORM connection
	UseGorm() *gorm.DB
}

type service struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	username   = os.Getenv("DB_USERNAME")
	password   = os.Getenv("DB_PASSWORD")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func StartDB() Service {

	// Keep the connection alive
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", username, password, host, port, database)
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// initalize GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Println("DB Ping Failed")
		log.Fatal(err)
	}

	dbInstance = &service{
		sqlDB:  sqlDB,
		gormDB: gormDB,
	}

	return dbInstance
}

func (s *service) UseGorm() *gorm.DB {
	return s.gormDB
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.sqlDB.Close()
}
