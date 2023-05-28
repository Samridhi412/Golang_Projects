package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/snowflakedb/gosnowflake"
)

func main() {
	var conn_string = "samridhi_garg1:Samridhi3404##@byjus-eds/TNL_TUTORPLUS/TNLTUTORPLUS?warehouse=SNOWFLAKE_WH_1&role=SNOWFLAKE_READER_1"
	db, err := sql.Open("snowflake", conn_string)
	if err != nil {
		log.Fatal(err)
	}else{
		fmt.Println("success")
	}
	// defer db.Close()
	var s sql.NullString
    id := 757882
	qry := "select count(*) from PUBLIC_ONE_TO_MANY_CLASSROOMS c "
	qry += "INNER JOIN PUBLIC_RAW_TOPICS rt on rt.ID=c.RAW_TOPIC_ID "
	qry += "INNER JOIN PUBLIC_CONTENT_BUNDLE_REQUISITES cbr on cbr.CONTENT_BUNDLE_ID=rt.CONTENT_BUNDLE_ID "
	qry += "INNER JOIN PUBLIC_ONE_TO_MANY_REQUISITE_GROUPS rg on rg.ID=cbr.REQUISITE_GROUP_ID "
	qry += "INNER JOIN PUBLIC_ONE_TO_MANY_REQUISITES r on r.REQUISITE_GROUP_ID=rg.ID where c.ID=" + fmt.Sprintf("'%d'", id)
	query := `CREATE TABLE IF NOT EXISTS packages (
		id int NOT NULL,
		name varchar NULL,
		created_at timestamp NOT NULL,
		modified_at timestamp NOT NULL,
		CONSTRAINT packages_pkey PRIMARY KEY (id)
	);`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)  
	defer cancelfunc()  
	res, err2 := db.ExecContext(ctx, query)  
	if err2 != nil {  
		log.Printf("Error %s when creating product table", err2)
		return
	}
	row, err3 := res.RowsAffected()  
    if err3 != nil {
        log.Printf("Error %s when getting rows affected", err3)
        return
    }
	log.Printf("Rows affected when creating table: %d", row)
	rows, err1 := db.Query(qry)
	if err1 != nil {
		log.Fatalf("Failed to scan. err: %v", err1)
		return
	}
	for rows.Next() {
		err := rows.Scan(&s)
		if err != nil {
			log.Fatalf("Failed to scan. err: %v", err)
		}
		if s.Valid {
			fmt.Println("Retrieved value:", s.String)
		} else {
			fmt.Println("Retrieved value: NULL")
		}
	}
}