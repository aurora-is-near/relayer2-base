package syncutils

import "sync"

type LockablePtr[T any] struct {
	ptr AtomicPtr[T]
	mtx sync.Mutex
}

func (lp *LockablePtr[T]) Load() *T {
	return lp.ptr.Load()
}

func (lp *LockablePtr[T]) LockIfNil() (value *T, unlock func(newValue *T)) {
	if value = lp.Load(); value == nil {
		lp.mtx.Lock()
		if value = lp.Load(); value == nil {
			return value, func(newValue *T) {
				lp.ptr.Store(newValue)
				lp.mtx.Unlock()
			}
		}
		lp.mtx.Unlock()
	}
	return value, nil
}

func (lp *LockablePtr[T]) LockIfNotNil() (value *T, unlock func(newValue *T)) {
	if value = lp.Load(); value != nil {
		lp.mtx.Lock()
		if value = lp.Load(); value != nil {
			return value, func(newValue *T) {
				lp.ptr.Store(newValue)
				lp.mtx.Unlock()
			}
		}
		lp.mtx.Unlock()
	}
	return value, nil
}
