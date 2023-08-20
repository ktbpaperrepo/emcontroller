package cpu

import (
	"fmt"
	"time"
)

var Workload int

// use CPU to do some calculation
func Exec() {
	multiSumCh()
}

func multiSumCh() {
	start := time.Now()
	chs := make([]chan int, 32)
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan int, 0)
		go sumCh(i, chs[i])
	}
	totalSum := 0
	for _, ch := range chs {
		sum := <-ch
		totalSum += sum
	}

	end := time.Now()
	fmt.Printf("total sum: %d, execution time: %g milliseconds\n", totalSum, float64(end.Sub(start).Microseconds())/1000)

}

func sumCh(seq int, ch chan int) {
	defer close(ch)
	sum := 0
	for i := 1; i <= Workload; i++ {
		sum += i
	}
	//fmt.Printf("goroutine %d result %d\n", seq, sum)
	ch <- sum
}
