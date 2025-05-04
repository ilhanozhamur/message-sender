package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ilhanozhamur/message-sender/internal/application"
)

type API struct {
	messageService *application.MessageService
	server         *http.Server
	toggleState    string
	mu             sync.Mutex
}

func NewAPI(service *application.MessageService) *API {
	return &API{
		messageService: service,
		toggleState:    "off",
	}
}

func (api *API) StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/start-stop", api.stateHandler)
	mux.HandleFunc("/sent-messages", api.sentMessagesHandler)

	api.server = &http.Server{Addr: ":8080", Handler: mux}
	api.server.ListenAndServe()
}

func (api *API) StopServer() {
	if api.server != nil {
		api.server.Close()
	}
}

func (api *API) stateHandler(w http.ResponseWriter, _ *http.Request) {
	if api.toggleState == "off" {
		api.toggleState = "on"
		go api.messageService.StartMessageSending()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sending started"))
	} else {
		api.toggleState = "off"
		go api.messageService.StopMessageSending()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sending stopped"))
	}
}

func (api *API) sentMessagesHandler(w http.ResponseWriter, _ *http.Request) {
	messages, err := api.messageService.GetSentMessages()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(messages) == 0 {
		http.Error(w, "No messages found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
