package api

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func InitDB() error {
	//connStr := os.Getenv("DB_CONNECTION_STRING")
	connStr := "database=promotions sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	db = conn

	err = db.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to the database")

	return nil
}

func storePromotion(promotion Promotion) error {
	stmt, err := db.Prepare("INSERT INTO promotions(id, price, expiration_date) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(promotion.ID, promotion.Price, promotion.ExpirationDate)
	if err != nil {
		return err
	}

	return nil
}

func getPromotion(id string) (Promotion, error) {
	query := `SELECT id, price, expiration_date FROM promotions WHERE id = $1`
	row := db.QueryRow(query, id)

	var promotion Promotion
	err := row.Scan(&promotion.ID, &promotion.Price, &promotion.ExpirationDate)
	if err != nil {
		log.Println(err)
		return Promotion{}, err
	}

	return promotion, nil
}

func TruncatePromotionsTable() error {
	_, err := db.Exec("TRUNCATE TABLE promotions")
	if err != nil {
		return err
	}

	log.Println("Table promotions truncated")
	return nil
}
