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

	// dynamic programming, l[i][j] = X means : substrings   s[i-X+1 : i+1] and t[j-X+1:j+1] are same
	// we comppute the values of l[i][j] and update the highest value of l[i][j] == longest substring == z

	for i := 0; i < r; i++ {
		for j := 0; j < n; j++ {
			if s[i] == t[j] {
				if i == 0 || j == 0 {
					l[i][j] = 1  // substrings of zero length are same
				} else {
					l[i][j] = l[i-1][j-1] + 1 // substrings s[i-z:i] === t[j-z:j] and s[i] == t[j] -> we add character s[i] and t[j] to substring
				}
				if l[i][j] > z { // if we found substring with biggest length
					z = l[i][j] // we update length of the new substring
					ret = s[i-z+1 : i+1] // and save the substring
				}
			} else { // s[i] != t[j] -> common substrings ending on i and j are of 0 length
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
