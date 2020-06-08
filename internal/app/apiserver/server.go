package apiserver

import (
	"GO-REST-API/internal/app/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	store store.Store
}

func newServer(store store.Store) *server{
	s := &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:store,
	}

	s.configureRouter()

	return s
}

func(s *server) configureRouter(){
	s.router.HandleFunc("/newrequest", s.newRequest()).Methods("POST")
}

func(s *server) ServeHTTP(w http.ResponseWriter,r *http.Request){
	s.router.ServeHTTP(w, r)
}

func(s *server) newRequest() http.HandlerFunc{

	return func(w http.ResponseWriter,r *http.Request){

	}
}