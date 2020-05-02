// Copyright 2020 @thiinbit(thiinbit@gmail.com). All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"sync"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

var (
	seqMu sync.Mutex
)

const (
	dbNamePrefix = "d_seq"
	tbNamePrefix = "t_seq"
	dbSize       = 4
	tbSize       = 8
)

func InitSeqDB() {
	initDB(dbNamePrefix, tbNamePrefix, dbSize, tbSize)
}
func DropSeqDB() {
	dropAllDB(dbNamePrefix, dbSize)
}

func initDB(dbNamePrefix string, tableNamePrefix string, dbSize int, tableSize int) {
	seqMu.Lock()
	defer seqMu.Unlock()

	for i := 0; i < dbSize; i++ {
		dbName := DBName(dbNamePrefix, i)

		createSql := "CREATE DATABASE IF NOT EXISTS " + dbName

		execSql(createSql)

		for j := 0; j < tableSize; j++ {
			tableName := TBName(dbName, tableNamePrefix, i*tableSize+j)

			createTableSql :=
				"CREATE TABLE IF NOT EXISTS " + tableName + " ( " +
					"`section` bigint(20) unsigned NOT NULL AUTO_INCREMENT, " +
					"`seq_step` bigint(20) unsigned NOT NULL, " +
					"`seq_max` bigint(20) unsigned NOT NULL, " +
					"`gmt_modify` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, " +
					"PRIMARY KEY (`section`) " +
					") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci; "

			execSql(createTableSql)
		}
	}
}

func dropAllDB(dbNamePrefix string, dbSize int) {
	seqMu.Lock()
	defer seqMu.Unlock()

	for i := 0; i < dbSize; i++ {
		dbName := DBName(dbNamePrefix, i)

		createSql := "DROP DATABASE " + dbName

		execSql(createSql)
	}
}

func execSql(sql string) {
	_, err := Exec(sql)
	if err != nil {
		panic(err)
	}
}
