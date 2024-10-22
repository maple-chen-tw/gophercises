package main

import (
	"database/sql"
	"fmt"
	"log"
)

func addNumberByRawSQL(db *sql.DB, number string) {
	sql := "INSERT INTO phone_numbers(number) VALUES(?)"
	result, err := db.Exec(sql, number)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return
	}

	//sql.Result 的LastInsertId()可取得AUTO_INCREMENT的值
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return
	}
	fmt.Println("Insert data id:", lastInsertID)

	//RowsAffected() 影響的資料筆數，如果很嚴謹的寫法會判斷RowsAffected()是否與新增的資料筆數一致
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)

}

func updateNumberByRawSQL(db *sql.DB, id string) {

	// 1. get current number
	var currentNumber string
	err := db.QueryRow("SELECT number FROM phone_numbers WHERE id=?", id).Scan(&currentNumber)
	if err != nil {
		fmt.Printf("Get number failed,err:%v\n", err)
		return
	}
	// 2. get normalization number
	newNumber := normalizedNumbers(currentNumber)

	// 3. check if exist same number
	var existingID string
	err = db.QueryRow("SELECT id FROM phone_numbers WHERE number=?", newNumber).Scan(&existingID)
	if err == nil {
		// if exist: delete the original entry.
		fmt.Printf("Get exist number %s with ID %s.", newNumber, existingID)
		_, err := db.Exec("DELETE FROM phone_numbers WHERE id=?", existingID)
		if err != nil {
			fmt.Printf("failed to delete exist number, err: %v\n", err)
			return
		}
	} else if err != sql.ErrNoRows {
		fmt.Printf("failed to check if exist same number, err:%v", err)
		return
	}

	// 4. update number
	sql := "UPDATE phone_numbers SET number=? WHERE id=?"
	result, err := db.Exec(sql, newNumber, id)
	if err != nil {
		fmt.Printf("Update data failed,err:%v", err)
		return
	}

	// 5.get affected row
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("Affected rows:", rowsaffected)
	fmt.Println("Update data successfully.")
}

func getAllIDs(db *sql.DB) ([]string, error) {
	var ids []string

	rows, err := db.Query("SELECT id FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
func createTable(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id INT AUTO_INCREMENT PRIMARY KEY,
		number VARCHAR(15) NOT NULL
	);`
	if _, err := db.Exec(sql); err != nil {
		log.Fatalf("create table failed: %v", err)
	}
	fmt.Println("create table successfully")
}
