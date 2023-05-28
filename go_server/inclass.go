package inclass
import (
	"database/sql"

	_ "github.com/snowflakedb/gosnowflake"
)

func main() {
	db, err := sql.Open("snowflake", "user:password@my_organization-my_account/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	...
}