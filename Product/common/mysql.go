package common

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// create MySql connection
func NewMysqlConn() (DB *sql.DB, err error) {
	DB, err = sql.Open("mysql", "root:04271017@tcp(127.0.0.1:3306)/product?charset=utf8")
	return DB, err
}

func GetResultRow(row *sql.Rows) map[string]string {
	column, _ := row.Columns()
	scanArgs := make([]interface{}, len(column))
	values := make([]interface{}, len(column))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]string)
	for row.Next() {
		if err := row.Scan(scanArgs...); err != nil {
			panic(err)
		}
		for i, col := range values {
			if col != nil {
				record[column[i]] = string(col.([]byte))
			}
		}
	}
	return record
}

func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	column, _ := rows.Columns()
	values := make([][]byte, len(column))
	scans := make([]interface{}, len(column))
	for k, _ := range values {
		scans[k] = &values[k]
	}
	i := 0
	results := make(map[int]map[string]string)
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			panic(err)
		}
		record := make(map[string]string)
		for k, v := range values {
			key := column[k]
			record[key] = string(v)
		}
		results[i] = record
		i++
	}
	return results
}
