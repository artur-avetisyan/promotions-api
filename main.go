package main

import (
	"csvParserAPI/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// please note that an overall project structure and naming should be improved as a part of production code
// just keeping the example app simple
func main() {
	err := api.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// this may not be the part of a business logic and is used here to test the example
	err = api.TruncatePromotionsTable()
	if err != nil {
		log.Fatal(err)
		return
	}

	// * although the parsing is making things safe, it may slow down things in production,
	// so we can consider removing actual parsing to float or time and handle db errors on insertion instead
	// keeping the parsing part for this example, may be removed in prod (so only insertion part is used)
	// ** keeping promotions.csv file inside project dir only to add parsing and storage logic as a part of this example
	err = api.ParseCSVAndStorePromotions("promotions.csv")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// most likely we will need an additional endpoint for receiving promotion csv files in production,
	// skipping that in this example, as the main logic is written above for parse and store and will be much similar
	router.HandleFunc("/promotions/{id}", api.GetPromotion).Methods("GET")

	log.Println("Server started on http://localhost:1321")
	log.Fatal(http.ListenAndServe(":1321", router))
}
