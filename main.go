package main

import (
	"brc/r6"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	cpuProfile := flag.String("profile", "", "write CPU profile to file")
	flag.Parse()
	f, err := os.Create(*cpuProfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	now := time.Now()
	inputFile := "measurements.txt"
	r6.R6(inputFile)
	fmt.Printf("took: %0.2f\n", time.Since(now).Seconds())
}
