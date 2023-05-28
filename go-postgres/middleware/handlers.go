package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct{
	ID int64 `json:"id,omitempty"`
	Message string `json:"message,omitempty`
}

func createConnection() *sql.DB{
    err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading .env file")
	}
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil{
		panic(err)
	}
    err = db.Ping()

	if err != nil{
		panic(err)
	}
	fmt.Println("Successfully connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request){
    var stock models.Stock
    err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil{
		log.Fatal("unable to descode the request body, %v",err)
	}
insertID := insertStock(stock)
res := response{
	ID: insertID,
	Message: "stock created successfully",
}
json.NewEncoder(w).Encode(res)

}

func GetStock(w http.ResponseWriter, r *http.Request){
    params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil{
      log.Fatal("Unable to convert string into int. %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil{
		log.Fatal("unable to get stock. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}

func GetAllStock(w http.ResponseWriter, r *http.Request){
   stocks, err := getAllStocks()
   if err != nil{
	log.Fatalf("unable to get all stocks %v", err)
   }
   json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request){
    params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil{
		log.Fatalf("unable to convert the string into int %v", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil{
		log.Fatalf("unable to descode the request body %v", err)
	}
	updatedRows := updateStock(int64(id), stock)
	msg := fmt.Sprintf("Stock updated successfully, Total rows affected %v", updatedRows)
	res:= response{
		ID: int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request){
 params := mux.Vars(r)
 id, err := strconv.Atoi(params["id"])
 if err != nil{
	log.Fatal("unable to convert string to int %v",err)
 }
 deletedRows := deleteStock(int64(id))
 msg := fmt.Sprintf("Stock deleted successfully, Total rows/records %v", deletedRows)
 res := response{
	ID: int64(id),
	Message: msg,
 }
 json.NewEncoder(w).Encode(res)
}
func insertStock(stock models.Stock) int64{
   db := createConnection()
   defer db.Close()
   sqlstatement := `INSERT INTO stocks(name, price, company) VALUES($1, $2, $3) RETURNING stockid`
   var id int64
   err := db.QueryRow(sqlstatement, stock.Name, stock.Price, stock.Company).Scan(&id)
   if err != nil{
	log.Fatalf("unable to execute query %v", err)
   }
   fmt.Printf("Inserted a single record %v", id)
   return id

}

func getStock(id int64)(models.Stock, error){
  db := createConnection()
  defer db.Close()
  var stock models.Stock
  sqlstatement := `SELECT * FROM stocks WHERE stockid=$1`
  row := db.QueryRow(sqlstatement, id)
  err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
  switch err {
  case sql.ErrNoRows:
	fmt.Println("No rows returned")
	return stock, nil
  case nil:
	return stock, nil
  default:
	log.Fatalf("unable to scan the row %v", err)		
  }
  return stock, err

}

func getAllStocks()([]models.Stock, error){
 db := createConnection()
 defer db.Close()
 var stocks []models.Stock
 sqlstatement := `SELECT * FROM stocks`
 rows, err := db.Query(sqlstatement)
 if err != nil{
	log.Fatalf("unable to execute the query %v", err)
 }
 defer rows.Close()
 for rows.Next(){
	var stock models.Stock
	err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	if err != nil{
		log.Fatalf("unable to scan the row %v", err)
	}
	stocks = append(stocks, stock)
 }
 return stocks, err
}

func updateStock(id int64, stock models.Stock) int64{
   db := createConnection()
   defer db.Close()
   sqlstatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
   res, err := db.Exec(sqlstatement, id, stock.Name, stock.Price, stock.Company) 
   if err != nil{
	log.Fatalf("unable to execute the query %v", err)
   } 
   rowsAffected, err := res.RowsAffected()
   if err != nil{
	log.Fatalf("Error while checking the affected rows %v", err)
   }
   fmt.Printf("Total rows/records affected %v", rowsAffected)
   return rowsAffected
}

func deleteStock(id int64) int64{
   db := createConnection()
   defer db.Close()
   sqlstatement := `DELETE FROM stocks WHERE stockid=$1`
   res, err := db.Exec(sqlstatement, id)
   if err != nil{
	log.Fatalf("unable to execute the query %v", err)
   }
   rowsaffected, err := res.RowsAffected()
   if err != nil{
	log.Fatalf("Error while checking the affected rows %v", err)
   }
   fmt.Printf("Total rows/records affected %v", rowsaffected)
   return rowsaffected
}