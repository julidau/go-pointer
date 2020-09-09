package pointer

// #include <stdlib.h>
import "C"
import (
	"sync"
	"unsafe"
)

var (
	mutex sync.RWMutex
	store = map[uintptr]interface{}{}
	count = uintptr(1) // start count at 1, since 0 is NULL and therefore special
)

func Save(v interface{}) (ptr unsafe.Pointer) {
	if v == nil {
		return nil
	}

	// Generate fake C pointer.
	// This pointer will not store any data, but will bi used for indexing purposes.
	mutex.Lock()
	ptr = unsafe.Pointer(count)
	store[count] = v
	count++
	mutex.Unlock()

	return ptr
}

func Restore(ptr unsafe.Pointer) (v interface{}) {
	if ptr == nil {
		return nil
	}

	mutex.RLock()
	v = store[uintptr(ptr)]
	mutex.RUnlock()
	return
}

func Unref(ptr unsafe.Pointer) {
	if ptr == nil {
		return
	}

	mutex.Lock()
	delete(store, uintptr(ptr))
	mutex.Unlock()
}
