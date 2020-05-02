// Copyright 2020 @thiinbit(thiinbit@gmail.com). All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	plog "github.com/thiinbit/plog4go"
)

const (
	seqDbLog       = "db-seq.log"
	driverName     = "mysql"
	dataSourceName = "root:root@/?parseTime=true&loc=Local"
)

// dbError is an error wrapper type for internal use only.
// Panics with errors are wrapped in dbError so that the top-level recover
// can distinguish intentional panics from this package.
type dbError struct{ error }

func handleRecover(de *dbError) {
	if r := recover(); r != nil {
		if dbErr, ok := r.(dbError); ok {
			de.error = dbErr.error
		} else {
			panic(r)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(dbError{err})
	}
}

func Exec(runSql string, args ...interface{}) (re sql.Result, err error) {
	defer handleRecover(&dbError{err})

	db, err := sql.Open(driverName, dataSourceName)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare(runSql)

	checkErr(err)

	re, err = stmt.Exec(args...)
	plog.Print("[SEQ_GEN] Run Exec SQL: ", runSql, args)
	checkErr(err)

	return re, err
}

func Query(runSql string, args ...interface{}) (rows *sql.Rows, err error) {
	defer handleRecover(&dbError{err})

	db, err := sql.Open(driverName, dataSourceName)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare(runSql)
	checkErr(err)

	rows, err = stmt.Query(args...)
	plog.Print("[SEQ_GEN] Run Query SQL: ", runSql, args)
	checkErr(err)

	defer stmt.Close()

	return rows, err
}

func QueryRow(runSql string, args ...interface{}) (row *sql.Row, err error) {
	defer handleRecover(&dbError{err})

	db, err := sql.Open(driverName, dataSourceName)
	checkErr(err)

	stmt, err := db.Prepare(runSql)
	checkErr(err)

	row = stmt.QueryRow(args...)
	plog.Print("[SEQ_GEN] Run Query SQL: ", runSql, args)
	checkErr(err)

	defer stmt.Close()

	return row, nil
}

func DBTBIndex(id int64, dbSize int, tbSize int) (int, int) {
	dbIdx, err := Mod(id, int64(dbSize))
	checkErr(err)
	tbIdxTemp, err := Mod(id, int64(tbSize))
	tbIdx := dbIdx*int64(tbSize) + tbIdxTemp
	checkErr(err)

	return int(dbIdx), int(tbIdx)
}

func DBTBName(id int64, dbPrefix string, tbPrefix string, dbSize int, tbSize int) (string, string) {
	dbIdx, tbIdx := DBTBIndex(id, dbSize, tbSize)

	dbName := DBName(dbPrefix, dbIdx)
	tbName := TBName(dbName, tbPrefix, tbIdx)

	return dbName, tbName
}

func DBName(dbPrefix string, idx int) string {
	return dbPrefix + fmt.Sprintf("%02d", idx)
}

// TBName return dbName.tableName
func TBName(dbName string, tbPrefix string, tbIdx int) string {
	return "`" + dbName + "`.`" + tbPrefix + fmt.Sprintf("%03d", tbIdx) + "`"
}
