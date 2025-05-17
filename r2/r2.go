package r2

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// 543 sec
func R2(inputFile string) {

	runtime.GOMAXPROCS(8)
	fmt.Println(runtime.NumCPU())

	var minMap, maxMap, sumMap, cntMap = make(map[string]float64), make(map[string]float64), make(map[string]float64), make(map[string]float64)
	type pair struct {
		station []byte
		temp    float64
	}

	var cMin, cMax, cSum = make(chan pair, 10_000), make(chan pair, 10_000), make(chan pair, 10_000)
	var conMin, conMax, conSum = true, true, true

	var fMin = func() {
		for p := range cMin {
			if cur, ok := minMap[string(p.station)]; ok {
				minMap[string(p.station)] = min(cur, p.temp)
			} else {
				minMap[string(p.station)] = p.temp
			}
		}
		conMin = false
	}
	var fMax = func() {
		for p := range cMax {
			if cur, ok := maxMap[string(p.station)]; ok {
				maxMap[string(p.station)] = max(cur, p.temp)
			} else {
				maxMap[string(p.station)] = p.temp
			}
		}
		conMax = false
	}
	var fSum = func() {
		for p := range cSum {
			cntMap[string(p.station)]++
			sumMap[string(p.station)] += p.temp
		}
		conSum = false
	}

	go fMin()
	go fMax()
	go fSum()

	file, _ := os.Open(inputFile)
	reader := bufio.NewReader(file)
	cnt := 0

	beginReadFile := time.Now()
	for {
		l, _, err := reader.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		station, temp := parseEntry(l)
		p := pair{station: station, temp: temp}
		cMin <- p
		cMax <- p
		cSum <- p
		cnt++
		if cnt%10_000_000 == 0 {
			fmt.Printf("processed: %d, took: %0.2f\n", cnt, time.Since(beginReadFile).Seconds())
		}
	}
	fmt.Printf("read file took: %02.f\n", time.Since(beginReadFile).Seconds())
	file.Close()

	close(cMin)
	close(cMax)
	close(cSum)
	for conMin || conMax || conSum {
	} // wait

	for k, _ := range minMap {
		println(k, minMap[k], maxMap[k], sumMap[k]/cntMap[k])
	}

}

func parseEntry(l []byte) (station []byte, temp float64) {
	s := string(l)
	sp := strings.Split(s, ";")
	station = []byte(sp[0])
	temp, _ = strconv.ParseFloat(sp[1], 64)
	return station, temp
}
