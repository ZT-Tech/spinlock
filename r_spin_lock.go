// Copyright 2019 darcy. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package spinlock

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// RTSpinLock re-entrant spin lock
type RTSpinLock struct {
	owner   int
	counter int
	state   int32
}

// NewRTSpinLock instantiates a spin-lock.
func NewRTSpinLock() sync.Locker {
	return new(RTSpinLock)
}

// Lock locks sl.
func (sl *RTSpinLock) Lock() {
	currentID := GetGoID()
	if sl.owner == currentID {
		sl.counter++
		return
	}
	for !sl.TryLock() {
		runtime.Gosched()
	}
	sl.owner = currentID
}

// Lock locks sl.
func (sl *RTSpinLock) Unlock() {
	if sl.owner != GetGoID() {
		fmt.Println("get goroutine id not eq")
		return
	}
	if sl.counter > 0 {
		sl.counter--
	} else {
		atomic.StoreInt32(&sl.state, 0)
	}
}

// TryLock try lock sl.
func (sl *RTSpinLock) TryLock() bool {
	return atomic.CompareAndSwapInt32(&sl.state, 0, mutexLocked)
}

// GetGoID get goroutine id
func GetGoID() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover:panic info:%v", err)
		}
	}()
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
