package api

import (
	"github.com/gorilla/mux"
	"net/http"
	log "github.com/sirupsen/logrus"
)

func NewRouter() *mux.Router{
	router := mux.NewRouter().StrictSlash(true)
	router.Use(loggingMiddleware)
	return router
}


func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Incoming %s from %s for %s",r.Method, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}