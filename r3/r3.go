package r3

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

func R3(inputFile string) {
	type stat struct {
		minm, maxx, summ float64
		cnt              int
	}

	var stationStat = make(map[string]*stat, 1_00_000)
	type pair struct {
		station string
		temp    float64
	}

	var cStat = make(chan pair, 100_000)

	file, _ := os.Open(inputFile)
	reader := bufio.NewReaderSize(file, 1<<19)
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
		// cStat <- p
		cur, ok := stationStat[p.station]
		if !ok {
			stationStat[p.station] = &stat{
				minm: p.temp,
				maxx: p.temp,
				summ: p.temp,
				cnt:  1,
			}
		} else {
			cur.cnt++
			cur.minm = min(cur.minm, p.temp)
			cur.maxx = max(cur.maxx, p.temp)
			cur.summ = min(cur.summ, p.temp)
			//stationStat[p.station] = cur
		}

		cnt++
		if cnt%10_000_000 == 0 {
			fmt.Printf("processed: %d, took: %0.2f\n", cnt, time.Since(beginReadFile).Seconds())
		}
	}
	fmt.Printf("read file took: %02.f\n", time.Since(beginReadFile).Seconds())
	file.Close()

	close(cStat)
	// wg.Wait()

	for k, v := range stationStat {
		println(k, v.minm, v.maxx, v.summ/float64(v.cnt))
	}

}

func parseEntry(l []byte) (station string, temp float64) {
	s := string(l)
	sp := strings.Split(s, ";")
	station = sp[0]
	temp, _ = strconv.ParseFloat(sp[1], 64)
	return station, temp
}
