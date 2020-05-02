// Copyright 2020 @thiinbit(thiinbit@gmail.com). All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"log"
	"testing"
	"time"
)

func TestNextSeq(t *testing.T) {
	start := time.Now()

	id, _ := Next()
	runTimes := 65536

	log.Print("id: ", id)

	var rv int64 = 9999
	for i := 0; i < runTimes; i++ {
		nv, err := NextSeq(id)
		if err != nil {
			t.Fatal(err)
		}

		if nv-1 != rv {
			t.Fatal("Wrong value")
		}

		rv = nv

	}

	elapsed := time.Since(start)
	log.Print("elapsed: ", elapsed)
	log.Print("last rv: ", rv)
}
