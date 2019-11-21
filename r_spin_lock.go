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
	sync.Mutex
	owner int64
	recursion int32
}

// NewRTSpinLock instantiates a spin-lock.
func NewRTSpinLock() sync.Locker {
	return new(RTSpinLock)
}

// Lock locks sl.
func (sl *RTSpinLock) Lock() {
	gid := GetGoID()
	if atomic.LoadInt64(&sl.owner) == gid {
		sl.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

// Lock locks sl.
func (sl *RTSpinLock) Unlock() {
	gid := GetGoID()
	if atomic.LoadInt64(&sl.owner) != gid {
		panic("unlock no permission")
	}
	m.recursion--
	if m.recursion != 0 {
		return 
	}
  	atomic.StoreInt64(&m.owner, -1)
  	m.Mutex.Unlock()
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
