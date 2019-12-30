package sufix

import (
	"encoding/json"
	"golang.org/x/net/html"
	"fmt"
	"net/http"
	"strings"
)

var num int
var delimeter int
var str string

var maxLen int

type OutEdge struct {
	anode      *Node
	labelStart int
	labelEnd   int
	bnode      *Node
}

type PK struct {
	node *Node
	key  string
}

type Node struct {
	parentkey    PK
	outedges     map[string]OutEdge
	suffixlink   *Node
	id           int
	outEdgesKeys []string
}

func NewNode(parentkey PK, outedges map[string]OutEdge, suffixlink *Node) Node {
	num += 1
	return Node{parentkey: parentkey, outedges: outedges, suffixlink: suffixlink, id: num - 1}
}

func (n *Node) Id() int {
	return n.id
}

func (n *Node) Suffixlink() *Node {
	return n.suffixlink
}

func (n *Node) SetSuffixlink(suffixlink *Node) {
	n.suffixlink = suffixlink
}

func (n *Node) Parentkey() PK {
	return n.parentkey
}

func (n *Node) SetParentkey(parentkey PK) {
	n.parentkey = parentkey
}

func (n *Node) setOutEdge(key string, anode *Node, startI int, endI int, bnode *Node) {
	//log.Printf("in setoutedge: key: %s, anodeid: %d, startI: %d, endI %d, bnode")
	//fmt.Printf("SETOUTEDGESTART node %d: key %s : A:id %d, start: %d, end: %d, bnode: %d\n",
	//	n.id, key, anode.id, startI, endI, bnode.id)

	if startI <= delimeter && endI > delimeter {
		endI = delimeter
	}

	if n.outedges == nil {
		n.outedges = make(map[string]OutEdge)
	}
	x := OutEdge{
		anode:      anode,
		labelStart: startI,
		labelEnd:   endI,
		bnode:      bnode,
	}
	n.outedges[key] = x
	n.outEdgesKeys = append(n.outEdgesKeys, key)
	//fmt.Printf("SETOUTEDGE node %d: key %s : A:id %d, start: %d, end: %d, bnode: %d\n",
	//	n.id, key, x.anode.id, x.labelStart, x.labelEnd, x.bnode.id)

}

func (n *Node) getOutEdge(key string) (anode *Node, labelStart int, labelEnd int, bnode *Node) {
	tmp := n.outedges[key]
	return tmp.anode, tmp.labelStart, tmp.labelEnd, tmp.bnode
}

func (n *Node) getOutEdges() map[string]OutEdge {
	return n.outedges
}

func in(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

func build(chars string) (Node, string) {
	pk := PK{}
	root := NewNode(pk, nil, nil)
	actnode := &root
	actkey := ""
	actlen, remainder, ind := 0, 0, 0

	for ind < len(chars) {
		//fmt.Printf("BUILDSTATE : nodeid: %d, actkey : %s, actlen %d, rem: %d, ind: %d \n",
		//	actnode.id, actkey, actlen, remainder, ind)
		ch := chars[ind : ind+1]
		if remainder == 0 {
			if len(actnode.outEdgesKeys) > 0 && in(actnode.outEdgesKeys, ch) {
				actkey = ch
				actlen = 1
				remainder = 1
				_, start, end, _ := actnode.getOutEdge(actkey)
				//fmt.Printf("EDGE : a: %d, start : %d, end %d, b: %d, \n", a.id, start, end, b.id)

				if end == maxLen {
					end = ind
				}
				if end-start+1 == actlen {
					_, _, _, actnode = actnode.getOutEdge(actkey)
					actkey = ""
					actlen = 0
				}
			} else {
				aleaf := NewNode(PK{}, nil, nil)
				pk := PK{
					node: actnode,
					key:  chars[ind : ind+1],
				}
				aleaf.SetParentkey(pk)
				//fmt.Printf("id: %d \n", actnode.id)
				actnode.setOutEdge(chars[ind:ind+1], actnode, ind, maxLen, &aleaf)
			}
		} else {
			if actkey == "" && actlen == 0 {
				if in(actnode.outEdgesKeys, ch) {
					actkey = ch
					actlen = 1
					remainder += 1
				} else {
					remainder += 1
					remainder, actnode, actkey, actlen =
						unfold(&root, chars, ind, remainder, actnode, actkey, actlen)
				}
			} else {
				_, start, end, _ := actnode.getOutEdge(actkey)
				if end == maxLen {
					end = ind
				}
				compareposition := start + actlen
				if chars[compareposition:compareposition+1] != ch {
					remainder += 1
					remainder, actnode, actkey, actlen =
						unfold(&root, chars, ind, remainder, actnode, actkey, actlen)
				} else {
					if compareposition < end {
						actlen += 1
						remainder += 1
					} else {
						remainder += 1
						_, _, _, tmp := actnode.getOutEdge(actkey)
						actnode = tmp
						if compareposition == end {
							actlen = 0
							actkey = ""
						} else {
							actlen = 1
							actkey = ch
						}
					}
				}
			}
		}
		ind += 1
	}
	return root, chars
}

func step(chars string, ind int, actnode *Node, actkey string, actlen int,
	remains string, ind_remainder int) (bool, *Node, string, int, int) {
	//fmt.Printf("in step, chars: %s, ind: %d, actnode: %d,"+
	//	" actkde: %s, actlen %d, remains: %s, ind_rem: %d \n",
	//	chars, ind, actnode.id, actkey, actlen, remains, ind_remainder)

	rem_label := remains[ind_remainder:]
	if actlen > 0 {
		_, start, end, _ := actnode.getOutEdge(actkey)
		if end == maxLen {
			end = ind
		}
		edgelabel := chars[start : end+1]
		if strings.HasPrefix(edgelabel, rem_label) {
			actlen = len(rem_label)
			actkey = rem_label[0:1]
			return true, actnode, actkey, actlen, ind_remainder
		}
	} else {
		if ind_remainder < len(remains) &&
			in(actnode.outEdgesKeys, remains[ind_remainder:ind_remainder+1]) {
			_, start, end, _ := actnode.getOutEdge(remains[ind_remainder : ind_remainder+1])
			if end == maxLen {
				end = ind
			}
			edgelabel := chars[start : end+1]
			if strings.HasPrefix(edgelabel, rem_label) {
				actlen = len(rem_label)
				actkey = rem_label[0:1]
				return true, actnode, actkey, actlen, ind_remainder
			}
		}
	}
	return false, actnode, actkey, actlen, ind_remainder
}

func hop(ind int, actnode *Node, actkey string, actlen int,
	remains string, ind_remainder int) (*Node, string, int, int) {
	//fmt.Printf("in hop, params: ind: %d, actnode: %d, actkey: %s, actlen %d, remains: %s, ind_remainder: %d\n", ind, actnode.id, actkey, actlen, remains, ind_remainder)
	if actlen == 0 || actkey == "" {
		//fmt.Printf("Leaving hop immediately\n")
		return actnode, actkey, actlen, ind_remainder
	}

	_, start, end, _ := actnode.getOutEdge(actkey)
	if end == maxLen {
		end = ind
	}
	edgelength := end - start + 1
	for actlen > edgelength {
		_, _, _, actnode = actnode.getOutEdge(actkey)
		ind_remainder += edgelength
		actkey = string(remains[ind_remainder : ind_remainder+1])
		actlen -= edgelength
		_, start, end, _ = actnode.getOutEdge(actkey)
		if end == maxLen {
			end = ind
		}
		edgelength = end - start + 1
	}
	if actlen == edgelength {
		_, _, _, actnode = actnode.getOutEdge(actkey)
		actkey = ""
		actlen = 0
		ind_remainder += edgelength
	}
	//fmt.Printf("LEAVIN Hop: actnode: %d, actkey: %s, actlen %d, ind_remainder: %d\n", actnode.id, actkey, actlen, ind_remainder)
	return actnode, actkey, actlen, ind_remainder
}

func unfold(root *Node, charsp string, indp int, remainderp int,
	actnode *Node, actkeyp string, actlenp int) (int, *Node, string, int) {
	var chars, ind, remainder, actkey, actlen = charsp, indp, remainderp, actkeyp, actlenp
	var prenode *Node = nil
	//var aleafLocB bool
	for remainder > 0 {
		aleaf := Node{}
		remains := chars[ind-remainder+1 : ind+1]
		actlen_re := len(remains) - 1 - actlen

		actnode, actkey, actlen, actlen_re = hop(ind, actnode, actkey, actlen, remains, actlen_re)
		var lost bool
		lost, actnode, actkey, actlen, actlen_re = step(chars, ind, actnode, actkey, actlen, remains, actlen_re)
		if lost {
			if actlen == 1 && prenode != nil && actnode.id != root.id {
				prenode.SetSuffixlink(actnode)
			}
			return remainder, actnode, actkey, actlen
		}
		//aleafLocB = false
		if actlen == 0 {
			if !in(actnode.outEdgesKeys, remains[actlen_re:actlen_re+1]) {
				aleaf = NewNode(PK{}, nil, nil)
				//aleafLocB = true
				//aedge :=  OutEdge{
				//	anode:      &actnode,
				//	labelStart: ind,
				//	labelEnd:   maxLen,
				//	bnode:      &aleaf,
				//}
				aleaf.SetParentkey(PK{
					node: actnode,
					key:  chars[ind : ind+1],
				})
				actnode.setOutEdge(chars[ind:ind+1], actnode, ind, maxLen, &aleaf)
			}
		} else {
			_, start, end, bnode := actnode.getOutEdge(actkey)
			if remains[actlen_re+actlen:actlen_re+actlen+1] != chars[start+actlen:start+actlen+1] {
				_, start, end, bnode = actnode.getOutEdge(actkey)
				//aleafLocB = true
				newnode := NewNode(PK{}, nil, nil)
				//halfedge1 := (actnode, start, start + actlen - 1, newnode)
				//halfedge2 := (newnode, start + actlen, end, bnode)
				actnode.setOutEdge(actkey, actnode, start, start+actlen-1, &newnode)
				newnode.SetParentkey(PK{actnode, actkey})
				newnode.setOutEdge(chars[start+actlen:start+actlen+1], &newnode, start+actlen, end, bnode)
				aleaf = NewNode(PK{}, nil, nil)
				//aedge := (newnode, ind, maxLen, aleaf)
				aleaf.SetParentkey(PK{
					node: &newnode,
					key:  chars[ind : ind+1],
				})
				newnode.setOutEdge(chars[ind:ind+1], &newnode, ind, maxLen, &aleaf)
			} else {
				return remainder, actnode, actkey, actlen
			}
		}

		if prenode != nil && aleaf.parentkey.node != nil {
			if aleaf.parentkey.node != root && prenode.id != aleaf.parentkey.node.id {
				prenode.SetSuffixlink(aleaf.parentkey.node)
			}
		}
		if aleaf.parentkey.node != nil && aleaf.parentkey.node != root {
			prenode = aleaf.Parentkey().node
		}
		if actnode.id == root.id && remainder > 1 {

			actkey = string(remains[1])
			actlen -= 1
		}
		if actnode.suffixlink != nil {
			actnode = actnode.suffixlink
		} else {
			actnode = root
		}
		remainder -= 1
	}
	return remainder, actnode, actkey, actlen
}

func printtree(n Node, chars string, depth int) {
	//fmt.Printf(" jsem node %d, suffix link: %d, ",
	//	root.id, root.suffixlink.id)
	fmt.Printf(strings.Repeat("--", depth))
	fmt.Printf("I am %d ", n.id)
	if n.suffixlink != nil {
		fmt.Printf("suffix link %d ", n.suffixlink.id)
	}
	fmt.Printf("\n")
	for x, val := range n.outedges {

		fmt.Printf(strings.Repeat("---", depth))
		fmt.Printf("edgekey: %s, edge: %d->%d, %s  %d..%d \n", x, val.anode.id, val.bnode.id,
			chars[val.labelStart:val.labelEnd+1], val.labelStart, val.labelEnd+1)
		printtree(*val.bnode, chars, depth+1)
	}
}

type Stack []*Node

func (s Stack) Push(v *Node) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, *Node) {
	// FIXME: What do we do if the stack is empty, though?

	l := len(s)
	return s[:l-1], s[l-1]
}

var stack Stack
var permaStack Stack
var flags map[int]bool
var nodeString map[int]string

func checkFlag(n *Node, o *OutEdge) {

}

//func LCSFromSuffixTreeNonRecursion(char string, n *Node, o *OutEdge, delimeter int, path string) (p string, flag int) {
//	var ret int
//	if len(n.outEdgesKeys) == 0 { // we are in leave
//		if o.labelStart <= delimeter && o.labelEnd <= delimeter+1 {
//			ret = -1
//		}
//		if o.labelStart >= delimeter && o.labelEnd > delimeter+1 {
//			ret = 1
//		}
//		p = ""
//		flag = ret
//		return
//	}
//
//	//maxLenFound := 0
//	//flags := make(map[int]bool)
//	//fupath := path
//
//	stackSize := 1
//	stack.Push(n)
//
//	// bfs using stack
//	for stackSize > 0 {
//		_, nd := stack.Pop()
//		stackSize -= 1
//		for _, val := range nd.outedges { // we go through deeper nodes
//			stack.Push(val.bnode)
//			stackSize += 1
//			permaStack.Push(val.bnode)
//		}
//	}
//
//}

func LCSFromSuffixTree(char string, n *Node, o *OutEdge, delimeter int, path string) (p string, flag int) {
	//fmt.Printf("\nin suff char %s, n.id %d, edge: %d..%d, del: %d, path: %s \n", char, n.id, o.labelStart,o.labelEnd+1, delimeter, path)

	// returns -1 if contains first string
	// returns 1 if contains second string
	// returns 0 if leaves contains both
	var ret int
	if len(n.outEdgesKeys) == 0 { // we are in leave
		if o.labelStart <= delimeter && o.labelEnd <= delimeter+1 {
			ret = -1
		}
		if o.labelStart >= delimeter && o.labelEnd > delimeter+1 {
			ret = 1
		}
		p = ""
		flag = ret
		return
	}

	maxLenFound := 0
	flags := make(map[int]bool)
	fupath := path

	for _, val := range n.outedges { // we go through deeper nodes
		if n.id != 0 {
			fupath = path + char[o.labelStart:o.labelEnd+1]
		}
		str, tmpflag := LCSFromSuffixTree(char, val.bnode, &val, delimeter, fupath)
		flags[tmpflag] = true // write found flags
		//fmt.Printf(" in id: %d, Writing flag: %b, of child : %d\n", n.id,tmpflag, val.bnode.id)
		if tmpflag == 0 && len(str) >= maxLenFound { // find the longest string with 0 flag
			p = str
			maxLenFound = len(str)
		}
	}
	//fmt.Printf("-j flag: %t, +1 flag : %t, 0 flag: %t \n", flags[-1], flags[+1], flags[0])

	if maxLenFound > 0 { // we found some string with
		flag = 0
		return
	} else { // we check if we are first zero
		if flags[-1] == true && flags[1] == true { // both substrings in leaves
			p = fupath
			flag = 0
			return
		} else {
			if flags[1] { //
				p = ""
				flag = 1
				return

			}
			if flags[-1] { // both substrings in leaves
				p = ""
				flag = -1
				return

			}
			panic("This should never happen")
		}
	}
}

//wrapper around this monstrosity
func LCSubstring(s string, t string) (ret string) {
	num, delimeter, str, maxLen = 0, 0, "", 0
	ending, delim := "$", "#"
	str = s + delim + t + ending
	delimeter = len(s)
	maxLen = len(str) - 1
	tree, _ := build(str)
	//printtree(tree, pst, 0)
	//stack = make([]*Node, maxLen)
	st, _ := LCSFromSuffixTree(str, &tree, &OutEdge{
		anode:      nil,
		labelStart: 0,
		labelEnd:   0,
		bnode:      nil,
	}, delimeter, "")
	return st

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
