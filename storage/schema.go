package storage

import (
	"github.com/nuveo/log"
)

func (s *Storage) CreateTables() {
	var tables = map[string]interface{}{
		"workflows": &Workflow{},
	}

	for k, v := range tables {
		s.db.DropTableIfExists(v)
		s.db.AutoMigrate(v)
		log.Printf("Table %s created\n", k)
	}
}

func (s *Storage) CreateSchema() {
	log.Println("Creating database schema")
	s.CreateTables()
}
