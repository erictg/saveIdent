package dbUtils



import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"log"
)


var db *sql.DB = nil


func GetDB() (*sql.DB, error){

	if db == nil{
		//todo TOML for connection strings
		connectionString := "root:rootpass@tcp(localhost:3306)/user_device_db"


		dbConn, err := sql.Open("mysql", connectionString)
		if err != nil{
			log.Println(err.Error())
			return nil, err
		}

		pingErr := dbConn.Ping()
		if pingErr != nil{
			log.Println(pingErr.Error())
			return nil, pingErr
		}


		//db is the package level db
		db = dbConn
	}

	return db, nil
}