// Copyright 2020 @thiinbit. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"errors"
	"time"
)

// 41bit time + 10bt machine + 11bit seq + 1bit chan = 63bit
// time use millids
const sequenceBits = 11

// Using two channels to average the probability of odd and even tail numbers
// and improve performance
const chanBits = 1
const machineIDBits = 10
const timeBits = 41

// 1
const sequenceShift = chanBits

// 12
const machineIDShift = sequenceBits + sequenceShift

// 22
const timeShift = machineIDBits + machineIDShift

// 00000000 00000000 00000000 00000000 00000000 00000000 00001111 11111110
const sequenceMask int64 = 4094

// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001
const chanMask int64 = 1

// TODO: Put into configuration later
const thisMachineID = 0

// timestamp of 2020,01,01 00:00:00:000 UTC
const sinceTime = 1577836800000

type generator struct {
	lastTimestamp int64
	sequence      int64
	genNum        int64
}

var genChan = make(chan *generator, 2)

// even tail generator
var gen0 = &generator{lastTimestamp: getTimestamp(), sequence: 0, genNum: 0}

// odd tail generator
var gen1 = &generator{lastTimestamp: getTimestamp(), sequence: 0, genNum: 1}

func init() {
	go func() {
		genChan <- gen0
		genChan <- gen1
	}()
}

// Next ID number
func Next() (int64, error) {
	_g := <-genChan
	defer func() { go func() { genChan <- _g }() }()

	if getTimestamp() < _g.lastTimestamp {
		return 0, errors.New("clock move backwards")
	}

	if getTimestamp() == _g.lastTimestamp {
		_g.sequence = (_g.sequence + (1 << 1)) & sequenceMask
		if _g.sequence>>1 == 0 {
			for getTimestamp() <= _g.lastTimestamp {
			}
		}
	} else {
		_g.sequence = 0
	}
	_g.lastTimestamp = getTimestamp()

	nextID := ((_g.lastTimestamp - sinceTime) << timeShift) | (thisMachineID << machineIDShift) | (_g.sequence << sequenceShift) | _g.genNum

	return nextID, nil
}

// ExtractTimestamp from ID
func ExtractTimestamp(seq int64) int64 {
	return sinceTime + (seq >> timeShift)
}

func getTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
