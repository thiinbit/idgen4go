// Copyright 2020 @thiinbit. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file

package idgen

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

// var lock sync.RWMutex

func init() {
	log.Println("Init max procs ", runtime.NumCPU()+1)
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
}

// TestUsage usage test
func TestUsage(t *testing.T) {
	// Generate a new ID
	id, _ := Next()

	// Get the ID generation time （int64 timestamp of millis）
	idTimestamp := ExtractTimestamp(id)

	fmt.Println(id, idTimestamp)
}

// TestNext
func TestNext(t *testing.T) {
	runTimes := 3
	testSize := 4000000
	printSize := 100
	checkDuplicate := false

	// var wg sync.WaitGroup
	// var lock sync.RWMutex

	counter := make(map[int64]int)

	timeCount := time.Duration(0)
	for rt := 0; rt < runTimes; rt++ {
		start := time.Now()

		for i := 0; i < testSize; i++ {

			// wg.Add(1)
			// go func(wg *sync.WaitGroup, idx int) {
			// defer wg.Done()

			id, err := Next()
			// _, err := Next()
			if err != nil {
				t.Fatal(err)
			}

			if checkDuplicate {
				// lock.Lock()
				counter[id]++
				// if idx < printSize {
				if i < printSize {
					t.Log(id)
				}
				// lock.Unlock()
			}

			// val, _ := counter.LoadOrStore(id, 0)
			// counter.Store(id, val.(int)+1)
			// }(&wg, i)
		}

		elapsed := time.Since(start)
		t.Log("Elapsed: ", elapsed)
		timeCount += elapsed
	}

	t.Log("AverageElapsed: ", timeCount.Seconds()/float64(runTimes))
	// wg.Wait()

	// counter.Range(func(k, v interface{}) bool {
	for k, v := range counter {
		// if v.(int) != 1 {
		if v != 1 {
			t.Fatalf("%d is appears %d times", k, v)
		}
	}
	// 	return true
	// })
}

// TestExtractTimestamp
func TestExtractTimestamp(t *testing.T) {
	testSize := 1000000

	start := time.Now()
	for i := 0; i < testSize; i++ {
		id, _ := Next()
		tT := ExtractTimestamp(id)
		if tT != getTimestamp() && tT != getTimestamp()-1 {
			t.Fatalf("id: %d, t: %d, extract timestampt fail!", id, tT)
		}
	}

	t.Log("Finished! Elapsed ", time.Since(start))
}

func TestExtractMachine(t *testing.T) {
	SetMachineID(128)
	n, _ := Next()

	if 128 != ExtractMachine(n) {
		t.Fatal("F")
	}

	SetMachineID(1024)
	n1, _ := Next()
	if 1024 == ExtractMachine(n1) || 0 != ExtractMachine(n1){
		log.Fatal("overflow")
	}

	for i := 0; i < 1024; i++ {
		SetMachineID(i)

		nextID, err := Next()
		if err != nil {
			t.Fatal(err)
		}

		mID := ExtractMachine(nextID)

		if i != mID {
			t.Fatalf("id %d, machine: %d, expect: %d", nextID, mID, i)
		}
	}

}
