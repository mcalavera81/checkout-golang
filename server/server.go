package main

import (
	"net/http"
	"checkout-service/server/api"
	"checkout-service/server/business"
	r "checkout-service/server/repository"
	"time"
	"checkout-service/server/config"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := api.NewRouter()

	api.ServeCheckoutResource(router, business.NewCheckoutService(r.NewCheckoutRepo(), r.NewCatalogRepo()))

	config := serverconfig.GetConfig()
	bindingAddr := fmt.Sprintf(":%d",config.Port)
	srv := &http.Server{
		Handler:      router,
		Addr:         bindingAddr,
		WriteTimeout: config.Timeout * time.Second,
		ReadTimeout:  config.Timeout * time.Second,
	}

	log.SetLevel(log.Level(config.LogLevel))
	log.Infof("Starting server at %s",bindingAddr)
	log.Fatal(srv.ListenAndServe())
}
