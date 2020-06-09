package apiserver

import (
	"GO-REST-API/internal/app/model"
	"GO-REST-API/internal/app/store"
	"encoding/json"
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
		req := &model.NewRequestRequest{}

		if err := json.NewDecoder(r.Body).Decode(req);err != nil{
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		if _,err := s.store.User().FindById(req.UserId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		rq := &model.NewRequestRequest{
			UserId: req.UserId,
			Message: req.Message,
		}

		if err := s.store.Request().NewRequest(rq); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w,r,http.StatusCreated, req)
	}

}

func(s *server) error(w http.ResponseWriter,r *http.Request, code int, err error){
	s.respond(w,r,code, map[string]string{"error":err.Error()})
}

func(s *server) respond(w http.ResponseWriter,r *http.Request, code int, data interface{}){
	w.WriteHeader(code)

	if data != nil{
		json.NewEncoder(w).Encode(data)
	}

}
