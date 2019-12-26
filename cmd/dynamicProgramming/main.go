package main

func main() {

}
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func LCSubstring(s string, t string) string {
	r, n, z := len(s), len(t), 0
	ret := ""
	l := make([][]int, r)
	for i := range l {
		l[i] = make([]int, n)
	}

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
					ret = s[i-z+1:i+1]
					println(ret)
				}
			}else {
				l[i][j] = 0
			}
		}
	}

	return ret

}

//
//function LCSubstr(S[1..r], T[1..n])
//L := array(1..r, 1..n)
//z := 0
//ret := {}
//
//for i := 1..r
//	for j := 1..n
//		if S[i] = T[j]
//			if i = 1 or j = 1
//				L[i, j] := 1
//			else
//				L[i, j] := L[i − 1, j − 1] + 1
//			if L[i, j] > z
//				z := L[i, j]
//				ret := {S[i − z + 1..z]}
//			//else if L[i, j] = z //
//				//ret := ret ∪ {S[i − z + 1..z]}
//		else
//			L[i, j] := 0
//return ret
