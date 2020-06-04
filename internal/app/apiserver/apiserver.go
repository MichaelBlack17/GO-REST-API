package apiserver

import (
	"GO-REST-API/internal/app/model"
	store2 "GO-REST-API/internal/app/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type  APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store *store2.Store
}

func New(config *Config) * APIServer{
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router:mux.NewRouter(),
	}
}
func (s *APIServer)Start() error{
	if err := s.configureLogger(); err!= nil{
		return err
	}

	s.configureRouter()

	if err := s.configureStore(); err != nil{
		return err
	}

	s.logger.Info("API starting...")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error  {
	level,err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil{
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureRouter(){
	s.router.HandleFunc("/hello", s.NewRequest())
}

func (s *APIServer) configureStore() error {
	st := store2.New(s.config.Store)
	if err := st.Open(); err != nil{
		return err
	}

	s.store = st
	return nil
}

func(s *APIServer) NewRequest() http.HandlerFunc{
	request := model.NewRequestRequest{}

	return func(w http.ResponseWriter, r *http.Request)	{
		req := &request
		if err := json.NewDecoder(r.Body).Decode(req); err != nil{
			s.error(w,r,http.StatusBadRequest, err)
			return
		}

		resp := model.NewRequestResponse{}
	}
}

func(s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error){
	s.respond(w, r, code, map[string]string{"error":err.Error()})
}

func(s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}){
	w.WriteHeader(code)
	if data != nil{
		json.NewEncoder(w).Encode(data)
	}
}