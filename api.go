package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type apiError struct {
	Error string
}

type APIServer struct {
	listenAddress string
	storage       Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		storage:       store,
	}
}

func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandleFunc(server.handleGetUser))
	router.HandleFunc("/user/{id}", makeHTTPHandleFunc(server.handleGetUserByID))
	router.HandleFunc("/user/create", makeHTTPHandleFunc(server.handleCreateUser))
	router.HandleFunc("/user/delete/{id}", makeHTTPHandleFunc(server.handleDeleteUser))
	router.HandleFunc("/user/update/{id}", makeHTTPHandleFunc(server.handleUpdateUser))

	log.Println("Server running on port: ", server.listenAddress)

	http.ListenAndServe(server.listenAddress, router)
}

// func (some poninter) someMethod(someParams, otherParams) someReturnType
// This structure binds the method to the "instance" of the pointer at the start
// Similar to "this.someMethod(someParams, otherParams)"
func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	users, err := s.storage.GetAllUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	user, err := s.storage.GetUserByID(idAsInt)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		return err
	}

	user := NewUser(createUserReq.Username, createUserReq.CurrentGame, createUserReq.CurrentLevel)
	if err := s.storage.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleUpdateUser(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	updateUserReq := new(User)
	if err := json.NewDecoder(r.Body).Decode(&updateUserReq); err != nil {
		return err
	}

	user := UpdateUser(updateUserReq.Username, updateUserReq.CurrentGame, updateUserReq.CurrentLevel, updateUserReq.CreatedAt)
	if err := s.storage.UpdateUser(idAsInt, user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	user, err := s.storage.DeleteUser(idAsInt)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
