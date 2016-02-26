package main

import (
	"runtime"
	"testing"
)

func Test_WriteInClosedChannel(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("Not recovered from writing to closed channel")
		}
	}()

	a := make(chan int)
	close(a)

	var b chan<- int
	b = a

	b <- 1

	t.Fatal("Should panic before on writing to closed channel")
}

func Test_CloseNullChannel(t *testing.T) {
	stage := 0

	defer func() {
		if x := recover(); x == nil {
			t.Fatal("Not recovered from closing closed channel")
		}

		if stage != 2 {
			t.Fatal("Invalid number of steps were done")
		}
	}()

	var c chan int

	c = nil
	// Can write to nil channel
	select {
	case c <- 1:
	default:
	}
	stage++
	// Can read from nil channel
	select {
	case <-c:
	default:
	}
	stage++
	// But can't close nil channel
	close(c)
	stage++

	t.Fatal("Should panic before on closing nil channel")
}

func Test_ReadFromBufClosedChannel(t *testing.T) {
	a := make(chan int, 1)
	a <- 1
	close(a)
	b := <-a
	if b != 1 {
		t.Fatal("Invalid b", b)
	}
}

func Test_LostWriteChannel(t *testing.T) {
	a := make(chan int)
	select {
	case a <- 1:
	default:
	}
	runtime.Gosched()
	select {
	case <-a:
		t.Fatal("Last channel write must be lost")
	default:
	}
}
