// Copyright 2019 darcy. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package spinlock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

const (
	mutexLocked = 1 << iota // mutex is locked
)

// NonRTSpinLock Non-reentrant spin lock
type NonRTSpinLock uint32

// NewNonRTSpinLock instantiates a spin-lock.
func NewNonRTSpinLock() sync.Locker {
	return new(NonRTSpinLock)
}

// Lock locks sl.
func (nrsl *NonRTSpinLock) Lock() {
	for !nrsl.TryLock() {
		runtime.Gosched()
	}
}

// Unlock unlocks sl.
func (nrsl *NonRTSpinLock) Unlock() {
	atomic.StoreUint32((*uint32)(nrsl), 0)
}

// TryLock try lock
func (nrsl *NonRTSpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32((*uint32)(nrsl), 0, mutexLocked)
}
