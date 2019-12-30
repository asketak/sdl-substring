package main

import (
	"fmt"
	"log"
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

	if startI<=delimeter  && endI >delimeter {
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
	fmt.Printf("SETOUTEDGE node %d: key %s : A:id %d, start: %d, end: %d, bnode: %d\n",
		n.id, key, x.anode.id, x.labelStart, x.labelEnd, x.bnode.id)

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
		fmt.Printf("BUILDSTATE : nodeid: %d, actkey : %s, actlen %d, rem: %d, ind: %d \n",
			actnode.id, actkey, actlen, remainder, ind)
		ch := chars[ind : ind+1]
		if remainder == 0 {
			if len(actnode.outEdgesKeys) > 0 && in(actnode.outEdgesKeys, ch) {
				println("1\n")
				actkey = ch
				actlen = 1
				remainder = 1
				a, start, end, b := actnode.getOutEdge(actkey)
				fmt.Printf("EDGE : a: %d, start : %d, end %d, b: %d, \n", a.id, start, end, b.id)

				if end == maxLen {
					println("2\n")
					end = ind
				}
				if end-start+1 == actlen {
					println("3\n")
					_, _, _, actnode = actnode.getOutEdge(actkey)
					actkey = ""
					actlen = 0
				}
			} else {
				println("4\n")
				aleaf := NewNode(PK{}, nil, nil)
				pk := PK{
					node: actnode,
					key:  chars[ind : ind+1],
				}
				aleaf.SetParentkey(pk)
				fmt.Printf("id: %d \n", actnode.id)
				actnode.setOutEdge(chars[ind:ind+1], actnode, ind, maxLen, &aleaf)
			}
		} else {
			println("a5\n")
			if actkey == "" && actlen == 0 {
				println("6\n")
				if in(actnode.outEdgesKeys, ch) {
					println("7\n")
					actkey = ch
					actlen = 1
					remainder += 1
				} else {
					log.Print("8\n")
					remainder += 1
					remainder, actnode, actkey, actlen =
						unfold(&root, chars, ind, remainder, actnode, actkey, actlen)
				}
			} else {
				println("9\n")
				a, start, end, b := actnode.getOutEdge(actkey)
				fmt.Printf(" a:id: %d,START: %d, end: %d b:id: %d,\n",
					a.id,start,end, b.id)
				if end == maxLen {
					println("10\n")
					end = ind
				}
				compareposition := start + actlen
				if chars[compareposition:compareposition+1] != ch {
					println("11\n")
					remainder += 1
					remainder, actnode, actkey, actlen =
						unfold(&root, chars, ind, remainder, actnode, actkey, actlen)
				} else {
					println("12\n")
					if compareposition < end {
						println("13\n")
						actlen += 1
						remainder += 1
					} else {
						println("14\n")
						remainder += 1
						_, _, _, tmp := actnode.getOutEdge(actkey)
						actnode = tmp
						if compareposition == end {
							println("1\n")
							actlen = 0
							actkey = ""
						} else {
							println("16\n")
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
	fmt.Printf("in step, chars: %s, ind: %d, actnode: %d,"+
		" actkde: %s, actlen %d, remains: %s, ind_rem: %d \n",
		chars, ind, actnode.id, actkey, actlen, remains, ind_remainder)

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
	fmt.Printf("in hop, params: ind: %d, actnode: %d, actkey: %s, actlen %d, remains: %s, ind_remainder: %d\n", ind, actnode.id, actkey, actlen, remains, ind_remainder)
	if actlen == 0 || actkey == "" {
		fmt.Printf("Leaving hop immediately\n")
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
	fmt.Printf("LEAVIN Hop: actnode: %d, actkey: %s, actlen %d, ind_remainder: %d\n", actnode.id, actkey, actlen, ind_remainder)
	return actnode, actkey, actlen, ind_remainder
}

func unfold(root *Node, charsp string, indp int, remainderp int,
	actnode *Node, actkeyp string, actlenp int) (int, *Node, string, int) {
	var chars, ind, remainder, actkey, actlen = charsp, indp, remainderp, actkeyp, actlenp
	println("IN UNFOLD")
	println("STATE")
	println(chars, ind, remainder, actkey, actlen)
	var prenode *Node = nil
	//var aleafLocB bool
	for remainder > 0 {
		aleaf := Node{}
		println("STATE")
		println(chars, ind, remainder, actkey, actlen)
		println("unfold remain: ", remainder)
		remains := chars[ind-remainder+1 : ind+1]
		println("KEYY: ", actkey)
		actlen_re := len(remains) - 1 - actlen

		actnode, actkey, actlen, actlen_re = hop(ind, actnode, actkey, actlen, remains, actlen_re)
		var lost bool
		lost, actnode, actkey, actlen, actlen_re = step(chars, ind, actnode, actkey, actlen, remains, actlen_re)
		if lost {
			println("unfold LOST")
			if actlen == 1 && prenode != nil && actnode.id != root.id {
				println("SETSUFIX  0")
				prenode.SetSuffixlink(actnode)
			}
			return remainder, actnode, actkey, actlen
		}
		//aleafLocB = false
		if actlen == 0 {
			println("unfold actlen=0")
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
			println("unfold actlen>0")
			_, start, end, bnode := actnode.getOutEdge(actkey)
			fmt.Printf("%s, %s, %d, %d, %d\n", remains, chars, start, actlen, actlen_re)
			fmt.Printf("%d\n", remainder)
			if actnode.suffixlink != nil {
				fmt.Printf("sufix: %d\n", actnode.suffixlink.id)
			}
			if remains[actlen_re+actlen:actlen_re+actlen+1] != chars[start+actlen:start+actlen+1] {
				println("unfold actlen chars")
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
				println("unfold return ", remainder, actnode.id, actkey, actlen)
				return remainder, actnode, actkey, actlen
			}
		}


		if prenode != nil && aleaf.parentkey.node != nil {
			fmt.Printf("prenode id: %d, aleafpk: %d", prenode.id, aleaf.parentkey.node.id)
			if aleaf.parentkey.node != root && prenode.id != aleaf.parentkey.node.id {
				println("SETSUFIX  1", prenode.id, "->", aleaf.parentkey.node.id)
				prenode.SetSuffixlink(aleaf.parentkey.node)
			}
		}
		if aleaf.parentkey.node != nil && aleaf.parentkey.node != root {
			prenode = aleaf.Parentkey().node
		}
		if actnode.id == root.id && remainder > 1 {
			println("ACK")
			println(remains)

			actkey = string(remains[1])
			actlen -= 1
			println(actkey, actlen)
		}
		if actnode.suffixlink != nil {
			println("SETSUFIX  2")
			actnode = actnode.suffixlink
		} else {
			actnode = root
		}
		remainder -= 1
		println("KEYY: ", actkey)
		println("ENDLOPPSTATE")
		println(remains, chars, ind, remainder, actkey, actlen)
	}
	println("unfold return ", remainder, actnode.id, actkey, actlen)
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

func LCSFromSuffixTree(char string, n Node, o OutEdge, delimeter int, path string) (p string, flag int) {
	fmt.Printf("\nin suff char %s, n.id %d, edge: %d..%d, del: %d, path: %s \n", char, n.id, o.labelStart,o.labelEnd+1, delimeter, path)

	// returns -1 if contains first string
	// returns 1 if contains second string
	// returns 0 if leaves contains both
	var ret int
	if len(n.outEdgesKeys) == 0 { // we are in leave
		println("INLEAVES")
		if o.labelStart <= delimeter && o.labelEnd <= delimeter {
			ret = -1
		}
		if o.labelStart >= delimeter && o.labelEnd >= delimeter {
			ret = 1
		}
		p = ""
		flag = ret
		println("leaving leave")
		println(p,flag)
		return
	}

	maxLenFound := 0
	flags := make(map[int]bool)
	fupath := path

	for _, val := range n.outedges { // we go through deeper nodes
		if n.id != 0{
			fupath = path + char[o.labelStart:o.labelEnd+1]
		}
		str, tmpflag := LCSFromSuffixTree(char, *val.bnode, val, delimeter, fupath)
		flags[tmpflag] = true                       // write found flags
		fmt.Printf(" in id: %d, Writing flag: %b, of child : %d\n", n.id,tmpflag, val.bnode.id)
		if tmpflag == 0 && len(str) >= maxLenFound { // find the longest string with 0 flag
			p = str
			maxLenFound = len(str)
		}
	}
	fmt.Printf("-j flag: %t, +1 flag : %t, 0 flag: %t \n", flags[-1], flags[+1], flags[0])


	if maxLenFound > 0 { // we found some string with
		flag = 0
		return
	} else { // we check if we are first zero
		if flags[-1] == true && flags[1] == true { // both substrings in leaves
			p = fupath
			flag = 0
			println("leave")
			println(p,flag)
			return
		} else {
			if flags[1] { //
				p = ""
				flag = 1
				println("leave")
				println(p,flag)
				return

			}
			if flags[-1] { // both substrings in leaves
				p = ""
				flag = -1
				println("leave")
				println(p,flag)
				return

			}
			panic("THIS HSOULD NOT HAPPEN")
		}
	}
}

//wrapper around this monstrosity
func LCSubstring(s string, t string) ( ret string )  {
	ending,delim := "$", "#"
	str = s + delim + t + ending
	delimeter = len(s)
	maxLen = len(str)-1
	tree, pst := build(str)
	println(tree.id)
	println(pst)
	fmt.Printf("'%s'", str)
	printtree(tree, pst, 0)
	st,_ := LCSFromSuffixTree(str, tree, OutEdge{
		anode:      nil,
		labelStart: 0,
		labelEnd:   0,
		bnode:      nil,
	}, delimeter,"")
	println("STRIIING")
	return st


}

func main() {
	//ending,delim := "$", "#"
	//str1 := strings.Repeat("na", 1000000)
	//str2 := strings.Repeat("ha", 1000000)

	//str1 := "xabxa"
	////str2 := "babxba"
	//str = str1 + delim+ str2 + ending
	//
	//
	//delimeter = len(str1)
	str := "aaa#aaaJJaaa$"
	delimeter = 3
	maxLen = len(str)-1
	tree, pst := build(str)
	println(tree.id)
	println(pst)
	fmt.Printf("'%s'", str)
	printtree(tree, pst, 0)
	//st,_ := LCSFromSuffixTree(str, tree, OutEdge{
	//	anode:      nil,
	//	labelStart: 0,
	//	labelEnd:   0,
	//	bnode:      nil,
	//}, delimeter,"")
	println("STRIIING")
	//println(st)
}
