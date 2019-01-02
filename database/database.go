package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bradfitz/slice"
	_ "github.com/lib/pq"
)

const source string = "postgres://%s:%s@%s/%s?sslmode=disable"

var db *sql.DB

func Connect(user, password, host, database string) {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf(source, user, password, host, database))
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("connected to: %s@%s\n", user, host)
}

func All(table string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", table)

	rows, err := db.Query(query)

	return processQuery(rows, err)
}

func Select(table string, column string, value string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", table, column)

	rows, err := db.Query(query, value)

	return processQuery(rows, err)
}

func Insert(table string, columns []string, values []interface{}) (int64, error) {
	vars := []string{}
	for i, _ := range values {
		vars = append(vars, fmt.Sprintf("$%s", strconv.Itoa(i+1)))
	}

	cols := strings.Join(columns, ", ")
	vals := strings.Join(vars, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, cols, vals)

	return executeQuery(query, values)
}

func Update(table string, column string, value string, columns []string, values []interface{}) (int64, error) {
	pairs := []string{}
	for i, n := range columns {
		pairs = append(pairs, fmt.Sprintf("%s = $%s", n, strconv.Itoa(i+2)))
	}

	fields := strings.Join(pairs, ", ")
	values = append([]interface{}{value}, values...)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $1", table, fields, column)

	return executeQuery(query, values)
}

func Delete(table string, column string, value string) (int64, error) {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", table, column)

	return executeQuery(query, []interface{}{value})
}

func RowsToList(rows *sql.Rows) ([]map[string]interface{}, error) {
	list := make([]map[string]interface{}, 0)
	cols, _ := rows.Columns()

	for rows.Next() {
		fields := make([]interface{}, len(cols))
		references := make([]interface{}, len(cols))
		for i, _ := range fields {
			references[i] = &fields[i]
		}

		if err := rows.Scan(references...); err != nil {
			return nil, err
		}

		item := make(map[string]interface{})
		for i, name := range cols {
			value := references[i].(*interface{})
			item[name] = *value
		}

		list = append(list, item)
	}

	return list, nil
}

func processQuery(rows *sql.Rows, err error) ([]map[string]interface{}, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list, err := RowsToList(rows)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v", list)
	fmt.Println()

	slice.Sort(list[:], func(x, y int) bool {
		return (list[x]["id"]).(int64) < (list[y]["id"]).(int64)
	})

	fmt.Printf("%v", list)
	fmt.Println()

	return list, nil
}

func executeQuery(query string, values []interface{}) (int64, error) {
	res, err := db.Exec(query, values...)
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return ra, nil
}
