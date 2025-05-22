package r6

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"time"
)

type stat struct {
	minm, maxx, cnt int32
	sum             int64
}
type partResult struct {
	no                           int
	discardFirstStringOfNextPart bool
	stationStat                  map[string]*stat
	firstString                  []byte // keep first string to handle overlapping data
}

func ReadPart(inputFile string, no, startByte, endByte int, result chan<- partResult) {
	stationStat := make(map[string]*stat, 10_000)

	file, _ := os.Open(inputFile)
	defer file.Close()
	file.Seek(int64(startByte), io.SeekStart)
	reader := bufio.NewReaderSize(file, 1<<19)
	readIdx := startByte
	cnt := 0
	var firstString []byte
	beginReadFile := time.Now()
	for readIdx < endByte {
		l, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		readIdx += len(l)

		if len(l) >= 1 {
			l = l[:len(l)-1]
		}
		if len(l) >= 1 && l[len(l)-1] == '\r' {
			l = l[:len(l)-1]
		}
		if cnt == 0 {
			firstString = l
			cnt++
			continue
		}
		cnt++
		if len(l) > 0 {
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
		}
		if cnt%1_000_000 == 0 {
			fmt.Printf("part %d processed: %d, took: %0.2f\n", no, cnt, time.Since(beginReadFile).Seconds())
		}
	}

	result <- partResult{
		no:                           no,
		stationStat:                  stationStat,
		firstString:                  firstString,
		discardFirstStringOfNextPart: readIdx > endByte,
	}
}

func R6(inputFile string) {
	stationStat := make(map[string]*stat, 10_000)
	file, _ := os.Open(inputFile)
	st, _ := file.Stat()
	sizeByte := int(st.Size())
	file.Close()
	parts := 25

	partResults := make(chan partResult, parts)

	for i := 0; i < parts; i++ {
		partBegin := sizeByte / parts * i
		partEnd := sizeByte/parts*(i+1) - 1
		go func(i int) {
			ReadPart(inputFile, i, partBegin, partEnd, partResults)
		}(i)
	}

	firsStringEachPart := map[int][]byte{}
	skippedParts := []int{}

	for _ = range parts {
		partRes := <-partResults
		firsStringEachPart[partRes.no] = partRes.firstString
		if partRes.discardFirstStringOfNextPart {
			skippedParts = append(skippedParts, partRes.no+1)
		}
		for station, stat := range partRes.stationStat {
			if stationStat[station] == nil {
				stationStat[station] = stat
			} else {
				stationStat[station].sum += stat.sum
				stationStat[station].cnt += stat.cnt
				stationStat[station].minm = min(stationStat[station].minm, stat.minm)
				stationStat[station].maxx = max(stationStat[station].maxx, stat.maxx)
			}
		}
	}
	for k, l := range firsStringEachPart {
		if slices.Contains(skippedParts, k) {
			continue
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
	}

	for k, v := range stationStat {
		fmt.Printf("%s %0.1f %0.1f %0.1f \n", k, float64(v.minm)/10.0, float64(v.maxx)/10.0, float64(v.sum)/10.0/float64(v.cnt))
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
