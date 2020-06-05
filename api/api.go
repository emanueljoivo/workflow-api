package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nuveo/log"
	"gitlab.com/emanueljoivo/workflows/messaging"
	"gitlab.com/emanueljoivo/workflows/storage"
	"net/http"
)

type HttpApi struct {
	Storage *storage.Storage
	Server  *http.Server
	Sender  *messaging.RabbitMQSender
	Consumer *messaging.RabbitMQConsumer
}

func NewAPI(storage *storage.Storage, sender *messaging.RabbitMQSender,
	consumer *messaging.RabbitMQConsumer) *HttpApi {
	return &HttpApi{
		Storage: storage,
		Sender:  sender,
		Consumer: consumer,
	}
}

func (a *HttpApi) Start(port string) error {
	a.Server = &http.Server{
		Addr:    ":" + port,
		Handler: a.bootRouter(),
	}
	log.Println("Service started")
	return a.Server.ListenAndServe()
}

func (a *HttpApi) Shutdown() error {
	return a.Server.Shutdown(context.Background())
}

func (a *HttpApi) bootRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/version", a.GetVersion).Methods(http.MethodGet)
	router.HandleFunc("/v1/workflow", a.CreateWorkflow).Methods(http.MethodPost)
	router.HandleFunc("/v1/workflow/{uuid}", a.UpdateWorkflowStatus).Methods(http.MethodPatch)
	router.HandleFunc("/v1/workflow", a.GetAllWorkflows).Methods(http.MethodGet)
	router.HandleFunc("/v1/workflow/consume", a.ConsumeWorkflow).Methods(http.MethodGet)

	return router
}

func Write(w http.ResponseWriter, statusCode int, i interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(i); err != nil {
		log.Println(EncodeResErr)
	}
}

type ErrorResponse struct {
	Message string `json:"Message"`
	Status  uint   `json:"Status"`
}