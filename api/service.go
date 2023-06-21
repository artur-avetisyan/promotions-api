package api

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetPromotion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	promotionID := vars["id"]

	promotion, err := getPromotion(promotionID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(promotion)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(w, string(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
}

func ParseCSVAndStorePromotions(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	reader := csv.NewReader(file)

	// Skip the header line if it exists
	//_, err = reader.Read()
	//if err != nil && err != io.EOF {
	//	return err
	//}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		id := record[0]
		priceStr := record[1]
		expirationDateStr := record[2]

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return err
		}

		expirationDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", expirationDateStr)
		if err != nil {
			return err
		}

		promotion := Promotion{
			ID:             id,
			Price:          price,
			ExpirationDate: expirationDate,
		}

		err = storePromotion(promotion)
		if err != nil {
			return err
		}
	}

	return nil
}
