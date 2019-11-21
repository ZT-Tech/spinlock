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

// TRSpinLock token re-entrant spin lock
type TRSpinLock struct {
	sync.Mutex
    token int64 
    recursion int32
}


// Lock locks sl.
func (m *TRMutex ) Lock(t int64)
{
  if atomic.LoadInt64(&m.token) == t{
    m.recursion++
    return
  }
  m.Mutex.Lock()
  atomic.StoreInt64(&m.token, t)
  m.recursion = 1
}

// Lock locks sl.
func (m *TRMutex ) Unlock(t int64)
{
  if atomic.LoadInt64(&m.token) != t{
    panic("token not eq")
  }
  m.recursion--
  if m.recursion != 0 {
    return 
  }
  atomic.StoreInt64(&m.token, 0)
  m.Mutex.Unlock()
}
