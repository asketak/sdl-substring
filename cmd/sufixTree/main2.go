package main

import (
	"fmt"
	"log"
	"strings"
)

var num int

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
	fmt.Printf("zapsano node %d: key %s : A:id %d, start: %d, end: %d, bnode: %d\n",
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
		print("nove chars:", chars)
		ch := chars[ind : ind+1]
		if remainder == 0 {
			if len(actnode.outEdgesKeys) > 0 && in(actnode.outEdgesKeys, ch) {
				print("1\n")
				actkey = ch
				actlen = 1
				remainder = 1
				_, start, end, _ := actnode.getOutEdge(actkey)
				if end == '#' {
					print("2\n")
					end = ind
				}
				if end-start+1 == actlen {
					print("3\n")
					_, _, _, actnode = actnode.getOutEdge(actkey)
					actkey = ""
					actlen = 0
				}
			} else {
				print("4\n")
				aleaf := NewNode(PK{}, nil, nil)
				pk := PK{
					node: actnode,
					key:  chars[ind : ind+1],
				}
				aleaf.SetParentkey(pk)
				println("pred setting")
				fmt.Printf("id: %d \n", actnode.id)
				actnode.setOutEdge(chars[ind:ind+1], actnode, ind, '#', &aleaf)
				for x, val := range actnode.outedges {
					fmt.Printf(" edgekey: %s, edge: %d->%d, str:%s..%s \n", x, val.anode.id, val.bnode.id,
						val.labelStart, val.labelEnd)
				}
			}
		} else {
			print("a5\n")
			if actkey == "" && actlen == 0 {
				print("6\n")
				if in(actnode.outEdgesKeys, ch) {
					print("7\n")
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
				print("9\n")
				_, start, end, _ := actnode.getOutEdge(actkey)
				if end == '#' {
					print("10\n")
					end = ind
				}
				compareposition := start + actlen
				if chars[compareposition:compareposition+1] != ch {
					print("11\n")
					remainder += 1
					remainder, actnode, actkey, actlen =
						unfold(&root, chars, ind, remainder, actnode, actkey, actlen)
				} else {
					print("12\n")
					if compareposition < end {
						print("13\n")
						actlen += 1
						remainder += 1
					} else {
						print("14\n")
						remainder += 1
						_, _, _, tmp := actnode.getOutEdge(actkey)
						actnode = tmp
						if compareposition == end {
							print("1\n")
							actlen = 0
							actkey = ""
						} else {
							print("16\n")
							actlen = 1
							actkey = ch
						}
					}
				}
			}
		}
		ind += 1
		if ind == len(chars) && remainder > 0 {
			print("pridavas/n")
			chars = chars + "$"
			print("nove chars:", chars)
		}
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
		if end == '#' {
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
			if end == '#' {
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
		fmt.Printf("Leaving hop immediately")
		return actnode, actkey, actlen, ind_remainder
	}

	_, start, end, _ := actnode.getOutEdge(actkey)
	if end == '#' {
		end = ind
	}
	edgelength := end - start + 1
	for actlen > edgelength {
		_, _, _, actnode = actnode.getOutEdge(actkey)
		ind_remainder += edgelength
		actkey = string(remains[ind_remainder : ind_remainder+1])
		actlen -= edgelength
		_, start, end, _ = actnode.getOutEdge(actkey)
		if end == '#' {
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
	print("IN UNFOLD")
	println("STATE")
	println(chars, ind, remainder, actkey, actlen)
	var prenode *Node = nil
	//var aleafLocB bool
	aleaf := Node{}
	for remainder > 0 {
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
			if actlen == 1 && prenode.id > 0 && actnode.id != root.id {
				println("SETSUFIX  0")
				prenode.SetSuffixlink(actnode)
			}
			return remainder, actnode, actkey, actlen
		}
		//aleafLocB = false
		if actlen == 0 {
			println("unfold actlen=0")
			if in(actnode.outEdgesKeys, remains[actlen_re:actlen_re+1]) {
				aleaf := NewNode(PK{}, nil, nil)
				//aleafLocB = true
				//aedge :=  OutEdge{
				//	anode:      &actnode,
				//	labelStart: ind,
				//	labelEnd:   '#',
				//	bnode:      &aleaf,
				//}
				aleaf.SetParentkey(PK{
					node: actnode,
					key:  chars[ind : ind+1],
				})
				actnode.setOutEdge(chars[ind:ind+1], actnode, ind, '#', &aleaf)
			}
		} else {
			println("unfold actlen>0")
			_, start, end, bnode := actnode.getOutEdge(actkey)
			fmt.Printf("%s, %s, %d, %d, %d\n", remains, chars, start, actlen, actlen_re)
			fmt.Printf("%d\n", remainder)
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
				//aedge := (newnode, ind, '#', aleaf)
				aleaf.SetParentkey(PK{
					node: &newnode,
					key:  chars[ind : ind+1],
				})
				newnode.setOutEdge(chars[ind:ind+1], &newnode, ind, '#', &aleaf)
			} else {
				println("unfold return ", remainder, actnode.id, actkey, actlen)
				return remainder, actnode, actkey, actlen
			}
		}
		if prenode != nil && aleaf.parentkey.node != nil {
			if aleaf.parentkey.node != root {
				println("SETSUFIX  1", prenode.id, "->", aleaf.parentkey.node.id)
				prenode.SetSuffixlink(aleaf.parentkey.node)
			}
		}
		if aleaf.parentkey.node != nil && aleaf.parentkey.node != root{
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

//func draw2(rnode Node, chars string, v int, ed string) {
//	l := len(chars)
//	edges := rnode.outEdgesKeys
//	nogc := make([]OutEdge, 0)
//	hasgc := make([]OutEdge, 0)
//	gc := make([]OutEdge, 0)
//	maxlen := len(chars) + 6
//	for _, edg := range rnode.getOutEdges() {
//		if v == 0 {
//			if len(edg.anode.suffixlink.outEdgesKeys) == 0 {
//				nogc = append(nogc, edg)
//			} else {
//				hasgc = append(hasgc, edg)
//			}
//		} else {
//			if len(edg.anode.suffixlink.outEdgesKeys) == 0 {
//				hasgc = append(nogc, edg)
//			} else {
//				nogc = append(hasgc, edg)
//			}
//		}
//	}
//	for _, x := range hasgc {
//		gc = append(gc, x)
//	}
//	for _, x := range nogc {
//		gc = append(gc, x)
//	}
//	for k, tmp := range gc{
//		parent, s,t, node := tmp.anode,tmp.labelStart,tmp.labelEnd, tmp.bnode
//
//	}
//
//
//}

//func draw1(root Node, chars string) {
//	fmt.Printf("\n %s \n (0)", chars)
//	draw2(root, chars, 0, "#")
//}

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

		fmt.Printf(strings.Repeat("--", depth))
		fmt.Printf("edgekey: %s, edge: %d->%d, %d..%d \n ", x, val.anode.id, val.bnode.id,
			val.labelStart, val.labelEnd)
		printtree(*val.bnode, chars, depth+1)
	}
}

func main() {
	str := "xabxaabxba"
	tree, pst := build(str)
	println(tree.id)
	println(pst)
	printtree(tree, pst, 0)
}
