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
	s.router.HandleFunc("/cancelrequest", s.cancelRequest()).Methods("DELETE")
	s.router.HandleFunc("/alluserrequests", s.allUserRequests()).Methods("GET")
	s.router.HandleFunc("/processingrequest", s.ProcessingRequest()).Methods("POST")
	s.router.HandleFunc("/cancelprocessingrequest", s.CancelProcessingRequest()).Methods("POST")
	s.router.HandleFunc("/allmanagerrequests", s.allManagerRequests()).Methods("GET")
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

func(s *server) cancelRequest() http.HandlerFunc{

	return func(w http.ResponseWriter,r *http.Request){
		req 	:= &model.CancelRequestRequest{}

		if err := json.NewDecoder(r.Body).Decode(req);err != nil{
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		if _,err := s.store.Request().FindByUserAndReqId(req.UserId, req.RequestId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		rq := &model.CancelRequestRequest{
			UserId: req.UserId,
			RequestId: req.RequestId,
		}

		resp, err := s.store.Request().CancelRequest(rq)

		if  err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w,r,http.StatusOK, resp)
	}

}

func(s *server) allUserRequests() http.HandlerFunc{

	return func(w http.ResponseWriter,r *http.Request){
		req := &model.AllUserRequestsRequest{}

		if err := json.NewDecoder(r.Body).Decode(req);err != nil{
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		if _,err := s.store.User().FindById(req.UserId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		rq := &model.AllUserRequestsRequest{
			UserId: req.UserId,
		}

		resp, err := s.store.Request().AllUserRequests(rq)
		if  err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w,r,http.StatusOK, resp)
	}

}

func(s *server) ProcessingRequest() http.HandlerFunc{

	return func(w http.ResponseWriter,r *http.Request){
		req := &model.ProcessingRequestRequest{}

		if err := json.NewDecoder(r.Body).Decode(req);err != nil{
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		if _,err := s.store.Manager().FindById(req.ManagerId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		if _,err := s.store.Manager().FindByManagerAndReqId(req.ManagerId, req.RequestId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		rq := &model.ProcessingRequestRequest{
			RequestId: req.RequestId,
			ManagerId: req.ManagerId,

		}

		resp, err := s.store.Request().ProcessingRequest(rq)
		if  err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w,r,http.StatusOK, resp)
	}

}

func(s *server) CancelProcessingRequest() http.HandlerFunc{

	return func(w http.ResponseWriter,r *http.Request){
		req := &model.CancelProcessingRequestRequest{}

		if err := json.NewDecoder(r.Body).Decode(req);err != nil{
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		if _,err := s.store.Manager().FindById(req.ManagerId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		if _,err := s.store.Manager().FindByManagerAndReqId(req.ManagerId, req.RequestId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		rq := &model.CancelProcessingRequestRequest{
			RequestId: req.RequestId,
			ManagerId: req.ManagerId,

		}

		resp, err := s.store.Request().CancelProcessingRequest(rq)
		if  err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w,r,http.StatusOK, resp)
	}

}

func(s *server) allManagerRequests() http.HandlerFunc{

	return func(w http.ResponseWriter,r *http.Request){
		req := &model.AllManagerRequestsRequest{}

		if err := json.NewDecoder(r.Body).Decode(req);err != nil{
			s.error(w,r, http.StatusBadRequest, err)
			return
		}

		if _,err := s.store.Manager().FindById(req.ManagerId); err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}

		resp, err := s.store.Request().AllManagerRequests(req)
		if  err != nil{
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w,r,http.StatusOK, resp)
	}

}

func(s *server) error(w http.ResponseWriter,r *http.Request, code int, err error){
	s.logger.Error(err)
	s.respond(w,r,code, map[string]string{"error":err.Error()})
}

func(s *server) respond(w http.ResponseWriter,r *http.Request, code int, data interface{}){

	w.WriteHeader(code)
	if data != nil{
		json.NewEncoder(w).Encode(data)
	}

}
