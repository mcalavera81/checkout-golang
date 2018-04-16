package api

import (
	"net/http"
	"checkout-service/server/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"strings"
	"checkout-service/server/business"
)

type(

	CheckoutResource struct {
		service business.CheckoutService
	}



)

func ServeCheckoutResource(router *mux.Router, service business.CheckoutService) {
	subrouter := router.PathPrefix("/api/checkout").Subrouter()
	resource := &CheckoutResource{service}
	subrouter.Methods("POST").Path("/").Name("Create").HandlerFunc(resource.Create)
	subrouter.Methods("PUT").Path("/{id:[a-zA-Z0-9]+}").Name("UpdateItem").HandlerFunc(resource.Update)
	subrouter.Methods("GET").Path("/{id:[a-zA-Z0-9]+}/total").Name("Total").HandlerFunc(resource.Total)
	subrouter.Methods("GET").Path("/{id:[a-zA-Z0-9]+}").Name("Get").HandlerFunc(resource.Get)
	subrouter.Methods("DELETE").Path("/{id:[a-zA-Z0-9]+}").Name("Get").HandlerFunc(resource.Delete)
}


func (resource *CheckoutResource) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketId, _ := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	basket, err := resource.service.Get(models.BasketId{basketId})

	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError, fmt.Sprintf("Error occured while querying basket %s ", basketId))
		return
	}

	if basket.IsEmpty() {
		errorHandler(w, models.NotFound, http.StatusNotFound, fmt.Sprintf("Basket %v not found", basketId))
		return
	}
	json.NewEncoder(w).Encode(&basket)
}


func (resource *CheckoutResource) Create(w http.ResponseWriter, r *http.Request) {
	basket, _ := resource.service.Create()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", strings.Join([]string{r.RequestURI, basket.Id},""))
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&basket)
}

func (resource *CheckoutResource) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketId := models.BasketId{vars["id"]}
	w.Header().Set("Content-Type", "application/json")


	item := new(models.CatalogItem)
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError, fmt.Sprintf("Error occured while updating basket %s", basketId))
	}

	basket, err := resource.service.Update(basketId, item)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError, fmt.Sprintf("Error occured while updating basket %s", basketId))
		return
	}

	json.NewEncoder(w).Encode(&basket)


}


func (resource *CheckoutResource) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketId, _ := vars["id"]
	w.Header().Set("Content-Type", "application/json")

	err := resource.service.Delete(models.BasketId{basketId})

	if err != nil {
		errorHandler(w, err, http.StatusNotFound, fmt.Sprintf("Basket %v not found", basketId))
		return
	}

}

func (resource *CheckoutResource) Total(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	basketId := models.BasketId{vars["id"]}
	w.Header().Set("Content-Type", "application/json")

	total, err := resource.service.Total(basketId)
	if err != nil {
		errorHandler(w, err, http.StatusInternalServerError, fmt.Sprintf("Error getting total for basket %s", basketId))
		return
	}

	json.NewEncoder(w).Encode(&struct {
		models.BasketId
		Total float64
	}{basketId, total})
}

