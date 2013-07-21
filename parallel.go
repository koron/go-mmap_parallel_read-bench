package main

import (
	"log"
)

const (
	BigFile = "big.txt"
)

func prepareBigFile() {
	// TODO:
	log.Println("prepareBigFile")
}

func sequentialRead() {
	// TODO:
	log.Println("sequentialRead")
}

func parallelRead() {
	// TODO:
	log.Println("parallelRead")
}

func warmup(f func()) {
	log.Println("warmup")
	for i := 0; i < 5; i++ {
		f()
	}
}

func benchmark(f func()) {
	// TODO:
	log.Println("benchmark")
	f()
}

func main() {
	prepareBigFile()
	warmup(sequentialRead)
	warmup(parallelRead)
	benchmark(sequentialRead)
	benchmark(parallelRead)
}
