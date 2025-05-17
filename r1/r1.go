package r1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	InputFile = "./measurements.txt"
)

// 300 sec
func R1() {

	var minMap, maxMap, sumMap, cntMap = make(map[string]float64), make(map[string]float64), make(map[string]float64), make(map[string]float64)

	file, _ := os.Open(InputFile)
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

		if _, ok := minMap[string(station)]; ok {
			minMap[string(station)] = min(minMap[string(station)], temp)
		} else {
			minMap[string(station)] = temp
		}
		if _, ok := maxMap[string(station)]; ok {
			maxMap[string(station)] = max(maxMap[string(station)], temp)
		} else {
			maxMap[string(station)] = temp
		}
		cntMap[string(station)]++
		sumMap[string(station)] += temp
		cnt++
		if cnt%10_000_000 == 0 {
			fmt.Printf("processed: %d, took: %0.2f\n", cnt, time.Since(beginReadFile).Seconds())
		}
	}

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
