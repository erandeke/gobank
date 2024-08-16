package routes

import (
	"encoding/json"
	"gobank/storage"
	"gobank/types"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	storage    storage.Store
}

func NewServer(listenAddr string, store storage.Store) *Server {
	return &Server{
		listenAddr: listenAddr,
		storage:    store,
	}
}

func (s *Server) routes() {
	routes := mux.NewRouter()

	routes.HandleFunc("/login", s.handleLogin).Methods("GET")
	routes.HandleFunc("/getaccount", makeHTTPHandleFunc(s.GetAccount)).Methods("GET")
	routes.HandleFunc("/createAccount", makeHTTPHandleFunc(s.CreateAccount)).Methods("POST")
}

func (s *Server) GetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)

}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) CreateAccount(w http.ResponseWriter, r *http.Request) error {

	req := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := types.NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	err = s.storage.CreateNewAccount(account)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)

}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(w)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, err)
		}
	}
}
