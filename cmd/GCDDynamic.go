package cloudSubstring

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

func LCSubstring(s string, t string) string {
	r, n, z := len(s), len(t), 0
	ret := ""
	l := make([][]int, r)
	for i := range l {
		l[i] = make([]int, n)
	}

	// dynamic programming, l[i][j] ==

	for i := 0; i < r; i++ {
		for j := 0; j < n; j++ {
			if s[i] == t[j] {
				if i == 0 || j == 0 {
					l[i][j] = 1
				} else {
					l[i][j] = l[i-1][j-1] + 1
				}
				if l[i][j] > z {
					z = l[i][j]
					ret = s[i-z+1 : i+1]
				}
			} else {
				l[i][j] = 0
			}
		}
	}

	return ret

}

// HelloHTTPMethod is an HTTP Cloud function.
// It uses the request method to differentiate the response.
func Hello(w http.ResponseWriter, r *http.Request) {
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

	response.Result = LCSubstring(request.S1, request.S2)

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
