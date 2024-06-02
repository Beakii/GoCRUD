package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// Creating a named type for the method signature used for API functions
type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//err := couldFail(param1, param2); err != nil
		//This patern scopes the err variable to the if so it doesnt get used further down
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type apiError struct {
	Error string
}

type APIServer struct {
	listenAddress string
}

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}

func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandleFunc(server.handleUser))

	log.Println("Server running on port: ", server.listenAddress)

	http.ListenAndServe(server.listenAddress, router)
}

// func (some poninter) someMethod(someParams, otherParams) someReturnType
// This structure binds the method to the "instance" of the pointer at the start
// Similar to "this.someMethod(someParams, otherParams)"
func (server *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return server.handleGetUser(w, r)
	case "POST":
		return server.handleCreateUser(w, r)
	case "PUT":
		return server.handleUpdateUser(w, r)
	case "DELETE":
		return server.handleDeleteUser(w, r)
	default:
		return fmt.Errorf("method not allowed  %s", r.Method)
	}
}

func (server *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	account := NewUser("Beakie", "Black Desert Online", 65)
	return WriteJSON(w, http.StatusOK, account)
}

func (server *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
