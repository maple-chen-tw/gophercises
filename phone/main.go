package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var phoneNumbers = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func main() {

	// set up database connection
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connect to [phone] Database！")

	// insert test data
	for _, number := range phoneNumbers {
		addNumberByRawSQL(db, number)
	}
	fmt.Println("Successfully add [Number] to [phone] Database！")

	// get all ids from database
	ids, err := getAllIDs(db)
	if err != nil {
		log.Fatalf("fail to get ID, err:%v\n", err)

	}
	//for _, id := range ids {
	//	fmt.Printf("ID: %s", id)
	//}

	// update data
	for _, id := range ids {
		updateNumberByRawSQL(db, id)
	}
	fmt.Println("Successfully update [Number] to [phone] Database！")

}

func normalizedNumbers(number string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(number, "")
}
