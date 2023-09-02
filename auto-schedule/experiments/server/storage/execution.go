package storage

import (
	"math"
	"os"
	"runtime/debug"
	"time"
)

var DiskCanUse float64 // unit GiB

const (
	// we do not occupy all accessible storage, because perhaps the program itself needs to use some.
	reserveStorage float64 = 2 // unit GiB
	occupyFileName string  = "occupy.tmp"
)

func Exec() {
	OccupyStorage()
}

// occupies the specific size (unit GiB) of storage
func OccupyStorage() {
	occupyManually()

	// manually release the memory occupied when writing file.
	time.Sleep(100 * time.Millisecond)
	debug.FreeOSMemory()
	time.Sleep(100 * time.Millisecond)
}

// the following 2 sparse file methods sometimes do not work, so I make this.
func occupyManually() {
	f, err := os.Create(occupyFileName)
	if err != nil {
		panic("occupyTruncate create file, error:" + err.Error())
	}
	defer f.Close()

	var contentMiB []byte = make([]byte, 1024*1024)                        // 1 MiB
	var numContent int = int(math.Floor(DiskCanUse-reserveStorage)) * 1024 // number of 1 MiB

	//var fileSize int64 = int64(math.Floor(DiskCanUse-reserveStorage)) * 1024 * 1024 * 1024 // unit B

	for i := 0; i < numContent; i++ {
		_, err = f.Write(contentMiB)
		if err != nil {
			panic("occupyManually Write file, error:" + err.Error())
		}
	}

}

// use truncate to make a sparse file to occupy storage, sometimes not work
func occupyTruncate() {
	f, err := os.Create(occupyFileName)
	if err != nil {
		panic("occupyTruncate create file, error:" + err.Error())
	}
	defer f.Close()

	var fileSize int64 = int64(math.Floor(DiskCanUse-reserveStorage)) * 1024 * 1024 * 1024 // unit B

	if err := f.Truncate(fileSize); err != nil {
		panic("occupyTruncate Truncate file, error:" + err.Error())
	}

	if err := f.Sync(); err != nil {
		panic("occupyTruncate Sync file, error:" + err.Error())
	}

}

// use seek to make a sparse file to occupy storage, sometimes not work
func occupySeek() {
	f, err := os.Create(occupyFileName)
	if err != nil {
		panic("occupySeek create file, error:" + err.Error())
	}
	defer f.Close()

	var fileSize int64 = int64(math.Floor(DiskCanUse-reserveStorage)) * 1024 * 1024 * 1024 // unit B

	_, err = f.Seek(fileSize-1, 0)
	if err != nil {
		panic("occupySeek Seek file, error:" + err.Error())
	}

	_, err = f.Write([]byte{0})
	if err != nil {
		panic("occupySeek Write file, error:" + err.Error())
	}

	if err := f.Sync(); err != nil {
		panic("occupySeek Sync file, error:" + err.Error())
	}

}
