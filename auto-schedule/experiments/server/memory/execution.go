package memory

import (
	"fmt"
	"math"
	"runtime"
	"runtime/debug"
	"time"
)

var MemCanUse float64 // unit MiB

// we do not occupy all accessible memory, because perhaps the program itself needs to use some.
const reserveMemory float64 = 1536 // unit MiB
const reserveMemPct float64 = 0.4

func Exec() {
	OccupyMem()
}

// occupies the specific size (unit MiB) of memory
func OccupyMem() {
	// an int64 variable occupies 8 byte, try fmt.Println(unsafe.Sizeof(int64(123)))
	// a byte variable occupies 1 byte, try fmt.Println(unsafe.Sizeof(byte('a')))

	// if we allocate too much memory at the same time, there may be some problems, so we should allocate Memory in chunks.

	chunkSize := 1024 * 1024 // we set 1 MiB is a chunk
	chunkNum := int(math.Floor((MemCanUse - reserveMemory) * (1 - reserveMemPct)))
	if chunkNum < 0 {
		chunkNum = 0 // if we do not do this, make([][]byte, chunkNum) will panic when chunkNum < 0
	}
	var memoryAllocated [][]byte = make([][]byte, chunkNum)

	for i := 0; i < chunkNum; i++ {
		memoryChunk := make([]byte, chunkSize) // the memory will not be released instantly
		memoryAllocated[i] = append(memoryAllocated[i], memoryChunk...)

		//memoryChunk = nil // I tried to instantly release the memory by this, but it does not work.

		// The GC of Go language is executed every 2 minutes, so the memory occupied by memoryChunk will not be released instantly,
		// so we manually call the GC when allocating every 100 MiB memory
		if i%100 == 0 {
			//runtime.GC()
			time.Sleep(100 * time.Millisecond)
			debug.FreeOSMemory() // sometimes runtime.GC() does not work. This is better.
			time.Sleep(100 * time.Millisecond)
		}

		// Sleep for a short period to avoid consuming too much CPU
		//time.Sleep(100 * time.Millisecond)
	}
	//runtime.GC() // GC after allocating all memory
	time.Sleep(100 * time.Millisecond)
	debug.FreeOSMemory() // sometimes runtime.GC() does not work. This is better.
	time.Sleep(100 * time.Millisecond)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated memory: %.2f MiB\n", float64(m.Alloc)/(1024*1024))

	// block forever
	tmpCh := make(chan int)
	select {
	case _ = <-tmpCh:
	}

	fmt.Println(memoryAllocated)
}
