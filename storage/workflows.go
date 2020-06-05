package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/nuveo/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueStr, err := json.Marshal(j)
	return string(valueStr), err
}

func (j *JSONB) Scan(value interface{}) error {
	src, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	var i interface{}
	if err := json.Unmarshal(src, &i); err != nil {
		return err
	}

	*j, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{} failed")
	}

	return nil
}

type WorkflowStatus uint8

const (
	Inserted WorkflowStatus = 0
	Consumed WorkflowStatus = 1
)

func (ws WorkflowStatus) String() string {
	return [...]string{"Inserted", "Consumed"}[ws]
}

type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (w *Base) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}

type Workflow struct {
	Base
	Status WorkflowStatus
	Data JSONB 	`gorm:"type:json;default '{}'"`
	Steps pq.StringArray `gorm:"type:varchar(100)[]"`
}

var (
	CreateWorkflowErr = errors.New("unable to create workflow")
)

func (s *Storage) SaveWorkflow(w *Workflow) (*Workflow, error) {
	// Transaction started to ensure that when retrieve the workflow return the same one that was entered
	tx := s.db.Begin()

	fetchedWorkflow := &Workflow{}

	if tx.Create(&w).Error != nil {
		tx.Rollback()
		log.Errorf("Rollback done because %s", CreateWorkflowErr.Error())
		return nil, errors.New(CreateWorkflowErr.Error())
	} else {
		tx.Last(fetchedWorkflow)
		tx.Commit()
	}

	return fetchedWorkflow, nil
}

func (s *Storage) GetWorkflow(uuid string) (*Workflow, error) {
	tx := s.db.Begin()
	defer tx.Commit()

	fetchedWorkflow := &Workflow{}

	tx.Where("id = ?", uuid).First(&fetchedWorkflow)

	if fetchedWorkflow.Data == nil {
		errMsg := "No workflow with specified id found"
		log.Errorf("%s\n", errMsg)
		return nil, errors.New(errMsg)
	}
	return fetchedWorkflow, nil
}

func (s *Storage) GetAllWorkflows() []*Workflow {
	tx := s.db.Begin()

	var workflows []*Workflow

	tx.Find(&workflows)

	tx.Commit()

	return workflows
}

func (s *Storage) UpdateWorkflowStatus(uuid string) *Workflow {
	tx := s.db.Begin()

	fetchedWorkflow := &Workflow{}

	tx.Where("id = ?", uuid).First(&fetchedWorkflow)

	fetchedWorkflow.Status = Consumed

	tx.Save(fetchedWorkflow)

	tx.Where("id = ?", uuid).First(&fetchedWorkflow)

	tx.Commit()

	return fetchedWorkflow
}
