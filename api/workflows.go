package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nuveo/log"
	"gitlab.com/emanueljoivo/workflows/storage"
	"net/http"
	"os"
)

type WorkflowReq struct {
	Data  storage.JSONB `json:"data"`
	Steps []string `json:"steps"`
}

type WorkflowResp struct {
	ID string
	Status string
	Data storage.JSONB
	Steps []string
}

func (a *HttpApi) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to create workflow received")
	var workflowReq WorkflowReq
	var workflow *storage.Workflow

	err := json.NewDecoder(r.Body).Decode(&workflowReq)

	if err != nil {
		log.Errorf("%s\n", err.Error())
		Write(w, http.StatusBadRequest, ErrorResponse{
			Message: "May the body is malformed. Check and try again.",
			Status: http.StatusBadRequest,
		})
	} else {
		workflow = &storage.Workflow{
			Status: storage.Inserted,
			Data:   workflowReq.Data,
			Steps:  workflowReq.Steps,
		}

		workflow, err = a.Storage.SaveWorkflow(workflow)

		if err != nil {
			Write(w, http.StatusInternalServerError, ErrorResponse{
				Message: "A internal error occur, try again later.",
				Status: http.StatusInternalServerError,
			})
		} else {
			a.Sender.Send(workflow)

			var workflowResp = WorkflowResp{
				ID:     workflow.ID.String(),
				Status: workflow.Status.String(),
				Data:   workflow.Data,
				Steps:  workflow.Steps,
			}

			Write(w, http.StatusCreated, workflowResp)
			log.Println("Request attempt")
		}
	}
}

type WorkflowNewStatus struct {
	Status string `json:"status"`
}

func (a *HttpApi) UpdateWorkflowStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to update the workflow status received")
	var reqBody WorkflowNewStatus
	var workflow *storage.Workflow

	params := mux.Vars(r)

	workflowID := params["uuid"]

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		log.Errorf("%s\n", err.Error())
		Write(w, http.StatusBadRequest, ErrorResponse{
			Message: "Not found valid UUID.",
			Status: http.StatusBadRequest,
		})
	} else {
		workflow, err = a.Storage.GetWorkflow(workflowID)

		if err != nil {
			Write(w, http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
				Status: http.StatusBadRequest,
			})
		} else {
			if workflow.Status.String() == storage.Consumed.String() {
				Write(w, http.StatusBadRequest, ErrorResponse{
					Message: "The workflow already done and has status Consumed.",
					Status: http.StatusBadRequest,
				})
			} else {
				workflow.Status = storage.Consumed
				workflow = a.Storage.UpdateWorkflowStatus(workflow.ID.String())

				var workflowResp = WorkflowResp{
					ID:     workflow.ID.String(),
					Status: workflow.Status.String(),
					Data:   workflow.Data,
					Steps:  workflow.Steps,
				}

				Write(w, http.StatusAccepted, workflowResp)
				log.Println("Request attempt")

			}
		}
	}
}

func (a *HttpApi) GetAllWorkflows(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to retrieve all workflows received")
	var allWorkflows []WorkflowResp

	workflows := a.Storage.GetAllWorkflows()
	
	for i := 0; i < len(workflows); i++ {
		current := workflows[i]

		var workflow = WorkflowResp{
			ID:     current.ID.String(),
			Status: current.Status.String(),
			Data:   current.Data,
			Steps:  current.Steps,
		}

		allWorkflows = append(allWorkflows, workflow)
		// maybe a pagination
	}

	if len(allWorkflows) == 0 {
		allWorkflows = make([]WorkflowResp, 0)
	}

	Write(w, http.StatusOK, allWorkflows)
	log.Println("Request attempt")
}

func (a *HttpApi) ConsumeWorkflow(w http.ResponseWriter, r *http.Request) {
	wfConsumed := a.Consumer.Consume()

	wf, err := a.Storage.GetWorkflow(wfConsumed.ID.String())

	if err != nil {
		log.Errorf("Workflow not found in db")
		Write(w, http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		})
	} else {
		wfConsumed = a.Storage.UpdateWorkflowStatus(wf.ID.String())

		GenerateCSVFile(wfConsumed)

		var wfResp = WorkflowResp{
			ID:     wfConsumed.ID.String(),
			Status: wfConsumed.Status.String(),
			Data:   wfConsumed.Data,
			Steps:  wfConsumed.Steps,
		}

		Write(w, http.StatusOK, wfResp)
		log.Println("Request attempt")
	}
}

type CSVField struct {
	Key string
	Value string
}

func GenerateCSVFile(workflow *storage.Workflow) {
	log.Println("Generating CSV file")
	var data []CSVField

	b, err := json.Marshal(workflow.Data)

	if err != nil {
		log.Errorf("%s\n", err.Error())
	} else {
		err = json.Unmarshal(b, &data)

		fileName := fmt.Sprintf("./data/%s.csv", workflow.ID.String())

		csvFile, err := os.Create(fileName)

		if err != nil {
			log.Errorf("%s\n", err.Error())
		}

		defer csvFile.Close()

		writer := csv.NewWriter(csvFile)

		for _, d := range data {
			var row []string
			row = append(row, d.Key)
			row = append(row, d.Value)
			err := writer.Write(row)
			if err != nil {
				log.Errorf("%s\n", err.Error())
			}
		}
		writer.Flush()
	}
}
