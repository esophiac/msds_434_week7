package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// define a structure for an item in the product catalog
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"IsAvailable"`
}

// GetProductHandler is used to get data inside the product defined on out product catalog
func GetProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Read JSON file
		data, err := ioutil.ReadFile("./data/data.json")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// write body with JSON data
		rw.Header().Add("content-type", "application/json")
		rw.WriteHeader(http.StatusFound)
		rw.Write(data)
	}
}

// Create a new product and add to our product score
func CreateProductHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		// read incoming JSON request body
		data, err := ioutil.ReadAll(r.Body)
		// return bad status request if there is no body
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// validate data is proper JSON
		var product Product
		err = json.Unmarshal(data, &product)
		if err != nil {
			rw.WriteHeader(http.StatusExpectationFailed)
			rw.Write([]byte("Invalid Data Format"))
			return
		}
		// append data to existing product list
		var products []Product
		data, err = ioutil.ReadFile("./data/data.json")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// load JSON file into array of products
		err = json.Unmarshal(data, &products)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Add new Product to list
		products = append(products, product)
		// write updated JSON file
		updatedData, err := json.Marshal(products)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = ioutil.WriteFile("./data/data.json", updatedData, os.ModePerm)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// return after writing body
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("Added New Product"))
	}
}

func main() {

	// create new router
	router := mux.NewRouter()

	// route to perspective handlers
	router.Handle("/", GetProductHandler()).Methods("GET")
	router.Handle("/", CreateProductHandler()).Methods("POST")

	// Create new server and assign the router
	server := http.Server{
		Handler: router,
		Addr:    ":9090",
	}

	fmt.Println("Starting Product Catalog server on Port 9090")
	// start server on port/host
	server.ListenAndServe()
}
