// Copyright 2020 @thiinbit(thiinbit@gmail.com). All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"fmt"
)

const (
	allColumn      = ""
	defaultSeqStep = 10000
)

func updateAndSelect(sectionId int64) (int64, error) {

	_, tbName := sDBTBName(sectionId)

	upsertSql := fmt.Sprintf("INSERT INTO %s (section, seq_step, seq_max) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE seq_max = seq_max + seq_step", tbName)
	_, err := Exec(upsertSql, sectionId, defaultSeqStep, defaultSeqStep)
	if err != nil {
		return 0, err
	}

	runSql := fmt.Sprintf("SELECT seq_max FROM %s WHERE section = ?", tbName)

	row, err := QueryRow(runSql, sectionId)
	if err != nil {
		return 0, err
	}

	var seqMax int64
	err = row.Scan(&seqMax)

	if err != nil {
		return 0, err
	}

	return seqMax, nil
}

func delete(sectionId int64) (int64, error) {
	_, tbName := sDBTBName(sectionId)

	runSql := fmt.Sprintf("DELETE FROM %s WHERE uid = ?", tbName)

	re, err := Exec(runSql, sectionId)
	if err != nil {
		return 0, err
	}

	return re.RowsAffected()
}

func sDBTBName(sectionId int64) (string, string) {
	return DBTBName(sectionId, dbNamePrefix, tbNamePrefix, dbSize, tbSize)
}
