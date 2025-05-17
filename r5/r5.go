package r5

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func R5(inputFile string) {
	type stat struct {
		minm, maxx, cnt int32
		sum             int64
	}

	var stationStat = make(map[string]*stat, 10_000)

	file, _ := os.Open(inputFile)
	defer file.Close()
	reader := bufio.NewReaderSize(file, 1<<19)
	cnt := 0

	beginReadFile := time.Now()
	for {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		station, temp := parseEntry(l)
		stationS := string(station)
		cur, ok := stationStat[stationS]
		if !ok {
			stationStat[stationS] = &stat{
				minm: temp,
				maxx: temp,
				sum:  int64(temp),
				cnt:  1,
			}
		} else {
			cur.cnt++
			cur.minm = min(cur.minm, temp)
			cur.maxx = max(cur.maxx, temp)
			cur.sum += int64(temp)
		}

		cnt++
		if cnt%100_000_000 == 0 {
			fmt.Printf("processed: %d, took: %0.2f\n", cnt, time.Since(beginReadFile).Seconds())
		}
	}
	fmt.Printf("read file took: %02.f\n", time.Since(beginReadFile).Seconds())

	for k, v := range stationStat {
		println(k, float64(v.minm)/10.0, float64(v.maxx)/10.0, float64(v.sum)/10.0/float64(v.cnt))
	}

}

func parseEntry(l []byte) ([]byte, int32) {
	var temp int32
	var station []byte

	end := len(l)
	switch {
	case l[end-4] == ';':
		station = l[0 : end-4]
		temp = int32(l[end-3]-'0')*10 + int32(l[end-1]-'0')
	case l[end-5] == ';' && l[end-4] == '-':
		station = l[0 : end-5]
		temp = -(int32(l[end-3]-'0')*10 + int32(l[end-1]-'0'))
	case l[end-5] == ';' && l[end-4] != '-':
		station = l[0 : end-5]
		temp = int32(l[end-4]-'0')*100 + int32(l[end-3]-'0')*10 + int32(l[end-1]-'0')
	case l[end-6] == ';' && l[end-5] == '-':
		station = l[0 : end-6]
		temp = -(int32(l[end-4]-'0')*100 + int32(l[end-3]-'0')*10 + int32(l[end-1]-'0'))
	default:
		panic(string(l))
	}

	return station, temp
}
