package storage

import (
	"testing"
)

func TestCreateTables(t *testing.T) {
	s := OpenDriver()
	defer CloseDriver(s, t)

	t.Run("assertion that there is no table initially", func(t *testing.T) {
		if s.db.HasTable(&Workflow{}) {
			t.Errorf("expected to have no table but has at least one")
		}
	})

	t.Run("assertion that all tables were created", func(t *testing.T) {
		s.CreateTables()
		if !s.db.HasTable(&Workflow{}) {
			t.Errorf("expected to have no table but has at least one")
		}
	})
}

func TestDropTables(t *testing.T) {
	s := OpenDriver()
	defer CloseDriver(s, t)

	t.Run("assert drop with tables", func(t *testing.T) {
		s.CreateTables()
		s.db.DropTableIfExists(&Workflow{})
		if s.db.HasTable(&Workflow{}) {

			t.Errorf("expected to have no table but has at least one")
		}
	})
}

func OpenDriver() *Storage {
	return NewStorage("127.0.0.1", "5432", "workflows-admin",
		"workflows-db", "postgres")
}

func CloseDriver(s *Storage, t *testing.T) {
	err := s.db.Close()

	if err != nil {
		t.Fail()
	}
}