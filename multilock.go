/*
Copyright 2017 Albert Tedja

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package multilock

import (
	"sort"
	"sync"
)

var locks = struct {
	sync.Mutex
	list map[string]chan byte
}{
	list: make(map[string]chan byte),
}

type AcquiredLock struct {
	chans []chan byte
}

// Attempts to lock multiple keys. Keys must be unique.
// Return an object that must be passed back to Unlock().
func Lock(locks ...string) AcquiredLock {
	locks = unique(locks)
	sort.Strings(locks)
	acqChans := make([]chan byte, 0, len(locks))

	// get the channels and attempt to acquire them
	for _, l := range locks {
		acqChans = append(acqChans, getChan(l))
	}
	for _, ch := range acqChans {
		<-ch
	}

	return AcquiredLock{
		chans: acqChans,
	}
}

// Unlocks an acquired lock. Must be called after Lock.
func Unlock(al AcquiredLock) {
	for _, ch := range al.chans {
		ch <- 1
	}
}

// Create and get the channel for the specified key.
func getChan(key string) chan byte {
	locks.Lock()
	defer locks.Unlock()
	if locks.list[key] == nil {
		locks.list[key] = make(chan byte, 1)
		locks.list[key] <- 1
	}
	return locks.list[key]
}

// Return a new string with unique elements.
func unique(arr []string) []string {
	if arr == nil || len(arr) <= 1 {
		return arr
	}

	found := map[string]bool{}
	result := make([]string, 0, len(arr))
	for _, v := range arr {
		if !found[v] {
			found[v] = true
			result = append(result, v)
		}
	}
	return result
}