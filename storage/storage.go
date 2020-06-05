package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nuveo/log"
)

type Storage struct {
	db *gorm.DB
}

const dbDialect string =  "postgres"
var DB *Storage

func NewStorage(host string, port string, user string, dbname string, password string) *Storage {
	log.Println("Creating database")

	dbConfig := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	db, err := gorm.Open(dbDialect, dbConfig)

	if err != nil {
		log.Fatal("No db where found")
	}

	if db != nil {
		err = db.DB().Ping()

		if err != nil {
			log.Fatal(err.Error())
		}

		DB = &Storage{
			db,
		}
	}

	log.Println("Database is up")
	return DB
}

func (s *Storage) Setup() {
	s.CreateSchema()
}

func (s *Storage) DB() *gorm.DB {
	return s.db
}
