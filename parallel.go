package main

import (
	"fmt"
	"github.com/koron/jvgrep/mmap"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	bigFile  = "big.txt"
	bigSize  = 100 * 1024 * 1024 // =100MB
	seed     = 123456
	unitSize = bigSize / 4
)

func randomData(size int64) []byte {
	b := make([]byte, size)
	r := rand.New(rand.NewSource(seed))
	for i, _ := range b {
		b[i] = byte(r.Intn(256))
	}
	return b
}

func prepareBigFile() error {
	// Check big data file existence.
	fi, err := os.Stat(bigFile)
	if err == nil {
		if fi.Size() == bigSize {
			return nil
		}
	}

	// Create a new big data file
	log.Println("creating big data file")
	return ioutil.WriteFile(bigFile, randomData(bigSize), 0644)
}

func countZero(b []byte) int {
	cnt := 0
	for _, d := range b {
		if d == 0 {
			cnt++
		}
	}
	return cnt
}

func sequentialRead() (int, error) {
	mf, err := mmap.OpenMemfile(bigFile)
	if err != nil {
		return 0, err
	}
	defer mf.Close()

	return countZero(mf.Data()), nil
}

func parallelRead() (int, error) {
	mf, err := mmap.OpenMemfile(bigFile)
	if err != nil {
		return 0, err
	}
	defer mf.Close()

	ch := make(chan int, 4)

	whole := mf.Data()
	for i := 0; i < 4; i++ {
		go func(b []byte) {
			ch <- countZero(b)
		}(whole[i*unitSize : (i+1)*unitSize-1])
	}

	cnt := 0
	for i := 0; i < 4; i++ {
		select {
		case c := <-ch:
			cnt += c;
		}
	}

	return cnt, nil
}

func warmup(f func() (int, error), s string) {
	log.Println("warming up:", s)
	for i := 0; i < 5; i++ {
		_, err := f()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func benchmark(f func() (int, error), s string) {
	log.Println("benchmark:", s)
	for i := 0; i < 10; i++ {
		log.Println("#", i)
		_, err := f()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// be not used now.
func fullBench() {
	warmup(sequentialRead, "sequential")
	warmup(parallelRead, "parallel")
	benchmark(sequentialRead, "sequential")
	benchmark(parallelRead, "parallel")
}

func simpleBenchmark(f func() (int, error), s string) {
	log.Println("simple benchmark:", s)
	t := time.Now()
	cnt, _ := f()
	d := time.Since(t)
	log.Println("   count=", cnt)
	log.Println("  duration=", d)
}

func main() {
	err := prepareBigFile()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 2 {
		if os.Args[1] == "sequential" {
			simpleBenchmark(sequentialRead, "sequential")
			return
		} else if os.Args[1] == "parallel" {
			simpleBenchmark(parallelRead, "parallel")
			return
		}
	}
	fmt.Println("USAGE: go run parallel.go {sequential|parallel}")
}
