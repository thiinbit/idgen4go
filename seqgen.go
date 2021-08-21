// Copyright 2020 @thiinbit(thiinbit@gmail.com). All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"strconv"
	"sync"
	"time"
)

var (
	sequenceIncrCache = cache.New(7*24*time.Hour, 24*time.Hour)
	sequenceMaxCache  = cache.New(7*24*time.Hour, 24*time.Hour)
	seqIncrMu         sync.Mutex
)

const (
	Decimal = 10
)

// NextSeqBySection return next sequence of the id
func NextSeq(id int64) (int64, error) {
	// TODO: id to section map. Temporary ID and SEC are one-to-one relationships

	return cacheIncrOrInit(id)
}

// cacheIncrOrInit increment if cache exist or init 1 when not exist.
func cacheIncrOrInit(id int64) (int64, error) {
	seqIncrMu.Lock()
	defer seqIncrMu.Unlock()

	idStr := strconv.FormatInt(id, Decimal)

	// 1. Get or init current max sequence
	maxSeq, err := getOrInitMaxSeqCache(id, idStr)
	if err != nil {
		return 0, err
	}

	var nextVal interface{}

	// 2. Incr and check update max sequence
	if cv, found := sequenceIncrCache.Get(idStr); found {
		if cv.(int64) > maxSeq {
			return 0, errors.New("exceed max sequence")
		}

		nextVal, err = sequenceIncrCache.IncrementInt64(idStr, 1)
		if err != nil {
			return 0, err
		}
	} else {
		nextVal = maxSeq
		sequenceIncrCache.SetDefault(idStr, maxSeq)
	}

	// 3. Check need update max sequence
	if nextVal.(int64) >= maxSeq {
		_, err = updateAndCacheMaxSeq(id, idStr)
		if err != nil {
			return 0, err
		}
	}

	return nextVal.(int64), nil
}

func getOrInitMaxSeqCache(id int64, idStr string) (int64, error) {
	var cMax interface{}
	var found bool
	if cMax, found = sequenceMaxCache.Get(idStr); !found {
		seqMax, err := updateAndCacheMaxSeq(id, idStr)
		if err != nil {
			return 0, err
		}
		cMax = seqMax
	}

	return cMax.(int64), nil
}

func updateAndCacheMaxSeq(id int64, idStr string) (int64, error) {
	seqMax, err := updateAndSelect(id)
	if err != nil {
		return 0, err
	}

	sequenceMaxCache.SetDefault(idStr, seqMax)

	return seqMax, nil
}

//func updateIfHalfSeqUsed(nextSeq int64, maxSeq int64) {
//	TODO:....
//}
