package async

import (
	"encoding/json"
	"fmt"
	"github.com/asketak/sdlSubstring/deploy/dynamic"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"sync"
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
	for i, _ := range s {
		for j, _ := range t {
			wg.Add(1)
			go LCSsubtask(s, t, i, j)
		}
	}
	wg.Wait()
	for i, value := range substringSizes {
		log.Print(i, "AAA", value)
	}
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

func Entry(w http.ResponseWriter, r *http.Request) {
	var request struct {
		S1 string `json:"s1"`
		S2 string `json:"s2"`
	}

	var response struct {
		Result string `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "500 - Failed to parse request", http.StatusInternalServerError)
		return
	}

	response.Result = dynamic.LCSubstring(request.S1, request.S2)

	j, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, "500 - Error crafting response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = fmt.Fprint(w, html.EscapeString(string(j)))
	if err != nil{
		http.Error(w, "500 - Error crafting response", http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusOK)

}
