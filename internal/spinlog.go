package internal

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

func (sl *spinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched()
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

// SpinLock creates a new spin-lock.
func SpinLock() sync.Locker {
	return new(spinLock)
}
