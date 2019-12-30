package LCSAsync

import (
	"sync"
	"time"
)

var substringSizes = make(map[int]map[int]int) // how long is common prefix if we start at index
var lock = &sync.RWMutex{}
var wg = &sync.WaitGroup{}

func LCSsubtask(s string, t string, i int, j int) {
	ln := 0
	for ; ln+i < len(s) && ln+j < len(t) && s[ln+i] == t[ln+j]; ln++ {
	}
	defer wg.Done()
	lock.Lock()
	defer lock.Unlock()
	substringSizes[i][j] = ln
}

func LCSAsync(s string, t string) string {
	maxlen, index := 0, 0

	for i, _ := range s {
		substringSizes[i] = make(map[int]int)
	}
	time.Sleep(5*time.Second)
	for i, _ := range s {
		for j, _ := range t {
			wg.Add(1)
			go LCSsubtask(s, t, i, j)
		}
	}
	wg.Wait()
	for i, _ := range s {
		for j, _ := range t {
			if substringSizes[i][j] > maxlen {
				index = i
				maxlen = substringSizes[i][j]
			}
		}
	}
	return s[index : index+maxlen]

}
