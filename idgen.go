// Copyright 2020 @thiinbit. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"errors"
	"sync"
	"time"
)

// 41bit time + 10bt machine + 11bit seq + 1bit chan = 63bit
// time use millis
const sequenceBits = 11

// Using two channels to average the probability of odd and even tail numbers
// and improve performance
const chanBits = 1
const machineIDBits = 10 //
const timeBits = 41

// 1
const sequenceShift = chanBits

// 12
const machineIDShift = sequenceBits + sequenceShift

// 22
const timeShift = machineIDBits + machineIDShift

// 00000000 00000000 00000000 00000000 00000000 00000000 00001111 11111110
const sequenceMask int64 = 0xFFE

// 00000000 00000000 00000000 00000000 00000000 00111111 11110000 00000000
const machineMask int64 = 0x3FF000

// 00000000 00000000 00000000 00000000 00000000 00000000 00000000 00000001
const chanMask int64 = 0x01

// TODO: Put into configuration later
var thisMachineID int64 = 0

// timestamp of 2020,01,01 00:00:00:000 UTC
const sinceTime = 1577836800000

var mu sync.Mutex

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
	genChan <- gen0
	genChan <- gen1
}

// SetMachineID (only need set once on startup)
func SetMachineID(mID int) {
	mu.Lock()
	defer mu.Unlock()

	thisMachineID = int64(mID)
}

// Next ID number
func Next() (int64, error) {
	_g := <-genChan
	defer func() { genChan <- _g }()

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

// ExtractMachineID from ID
func ExtractMachine(seq int64) int {
	return int((seq & machineMask) >> machineIDShift)
}

// Hash
func Mod(id int64, m int64) int64 {
	return ExtractTimestamp(id) % m
}

func getTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
