package storage

import (
	"fmt"
	"testing"
)

func TestSaveWorkflow(t *testing.T) {
	s := OpenDriver()
	s.CreateTables()

	t.Run("assert that the insertion works", func (t *testing.T){
		if !s.db.HasTable(&Workflow{}) {
			t.Errorf("expected that to have a table created but not found")
		}

		var fakeData JSONB
		var fakeSteps []string

		fakeData = map[string]interface{}{
			"some": "data",
			"golang": "best lang",
			"number": 1,
		}

		var i int
		for i = 0; i < 100; i++ {
			fakeSteps = append(fakeSteps, fmt.Sprintf("step %d", i))
		}

		workflow := &Workflow{
			Status: WorkflowStatus(0),
			Data: fakeData,
			Steps:  fakeSteps,
		}

		insertedWorkflow, err := s.SaveWorkflow(workflow)

		if err != nil {
			t.Error(insertedWorkflow)
		}
	})

	s.db.Close()
}