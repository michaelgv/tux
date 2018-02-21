package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"log"
)

type Database struct {
	db *sql.DB
}

type Rows struct {
	rows *sql.Rows
}

type Row struct {
	row *sql.Row
}

type Result struct {
	result sql.Result
}

func GetDatabaseString() string {
	s := "root@/tux"
	return s
}

func MakeDatabase() *Database {
	this := new(Database)
	db, err := sql.Open("mysql", GetDatabaseString())
	checkErr(err)
	this.db = db
	return this
}

func (this *Database) Query(q string, args ...interface{}) Rows {
	log.Printf("%s on %v", q, args)
	rows, err := this.db.Query(q, args...)
	checkErr(err)
	return Rows{rows}
}

func (this *Database) QueryRow(q string, args ...interface{}) Row {
	log.Printf("%s on %v", q, args)
	row := this.db.QueryRow(q, args...)
	return Row{row}
}

func (this *Database) Exec(q string, args ...interface{}) Result {
	log.Printf("%s on %v", q, args)
	result, err := this.db.Exec(q, args...)
	checkErr(err)
	return Result{result}
}

func (r Rows) Close() {
	err := r.rows.Close()
	checkErr(err)
}

func (r Rows) Next() bool {
	return r.rows.Next()
}

func (r Rows) Scan(dest ...interface{}) {
	err := r.rows.Scan(dest...)
	checkErr(err)
}

func (r Row) Scan(dest ...interface{}) {
	err := r.row.Scan(dest...)
	checkErr(err)
}

func (r Result) LastInsertId() int {
	id, err := r.result.LastInsertId()
	checkErr(err)
	return int(id)
}

func (r Result) RowsAffected() int {
	count, err := r.result.RowsAffected()
	checkErr(err)
	return int(count)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}