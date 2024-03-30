package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBClient struct {
	db *sql.DB
}

func NewDBClient(connection string) (*DBClient, error) {
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DBClient{db}, nil
}

func (client *DBClient) GetTables() ([]string, error) {
	rows, err := client.db.Query("SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("Failed to get tables from database: %w", err)
	}

	defer rows.Close()

	var tableNames []string

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("Failed to scan table name: %w", err)
		}
		tableNames = append(tableNames, table)
	}

	return tableNames, nil
}

func (client *DBClient) GetRecords(table string, where string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)

	if where != "" {
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s", table, where)
	}

	rows, err := client.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get records from table: %w", err)
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Failed to get columns: %w", err)
	}

	var records []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("Failed to scan row: %w", err)
		}

		record := make(map[string]interface{})

		for i, col := range columns {
			val := values[i]

			switch val.(type) {
			case []byte:
				record[col] = string(val.([]byte))
			default:
				record[col] = val
			}
		}

		records = append(records, record)
	}

	return records, nil
}

func (client *DBClient) GetColumns(tableName string) (results []string, err error) {
	rows, err := client.db.Query("DESCRIBE " + tableName)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var columns []string

	for rows.Next() {
		var (
			column     string
			dataType   string
			null       string
			key        string
			defaultVal sql.NullString
			extra      string
		)

		if err := rows.Scan(&column, &dataType, &null, &key, &defaultVal, &extra); err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}
