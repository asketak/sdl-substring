package main

const oo = 1024 * 1024 * 32
const ALPHABET_SIZE = 256
const MAXN = 5000

var (
	Root       int
	Last_added int
	Pos        int
	NeedSl     int
	Remainder  int
	ActiveNode int
	ActiveEdge int
	ActiveLen  int
)

var tree [2 * MAXN]Node
var text [MAXN]byte

type Node struct {
	start int
	end   int
	slink int
	next  [ALPHABET_SIZE]int
}

func (n Node) edgeLength() int {
	tmp := n.end
	if Pos+1 < tmp {
		tmp = Pos + 1
	}
	return tmp - n.start
}

func new_node(start int, end int) int {
	var nd Node
	nd.start = start
	nd.end = end
	nd.slink = 0
	//for i:=0; i < ALPHABET_SIZE; i++ {
	//	nd.next[i] = 0
	//}
	Last_added += 1
	tree[Last_added] = nd
	return Last_added
}

func activeEdge() byte {
	return text[ActiveEdge]
}

func addSL(node int) {
	if NeedSl > 0 {
		tree[NeedSl].slink = node
	}
	NeedSl = node
}

func walkDown(node int) bool {
	if ActiveLen >= tree[node].edgeLength() {
		ActiveEdge += tree[node].edgeLength()
		ActiveLen -= tree[node].edgeLength()
		ActiveNode = node
		return true
	}
	return false
}
func st_init() {
	NeedSl, Last_added, Pos = 0, 0, -1
	Remainder, ActiveNode, ActiveEdge, ActiveLen = 0, 0, 0, 0
	ActiveNode = new_node(-1, -1)
	Root = ActiveNode
}

func st_extend(c byte) {
	Pos += 1
	text[Pos] = c
	NeedSl = 0
	Remainder += 1
	for Remainder > 0 {
		if ActiveLen == 0 {
			ActiveEdge = Pos
		}
		if tree[ActiveNode].next[activeEdge()] == 0 {
			leaf := new_node(Pos, oo)
			tree[ActiveNode].next[activeEdge()] = leaf
			addSL(ActiveNode) //rule 2
		} else {
			nxt := tree[ActiveNode].next[activeEdge()]
			if walkDown(nxt) {
				continue
			} //observation 2
			if text[tree[nxt].start+ActiveLen] == c { //observation 1
				ActiveLen += 1
				addSL(ActiveNode) //observation 3
				break
			}
			split := new_node(tree[nxt].start, tree[nxt].start+ActiveLen)
			tree[ActiveNode].next[activeEdge()] = split
			leaf := new_node(Pos, oo)
			tree[split].next[c] = leaf
			tree[nxt].start += ActiveLen
			tree[split].next[text[tree[nxt].start]] = nxt
			addSL(split) //rule 2
		}
		Remainder -= 1
		if ActiveNode == Root && ActiveLen > 0 { //rule 1
			ActiveLen -= 1
			ActiveEdge = Pos - Remainder + 1
		} else {
			if tree[ActiveNode].slink > 0 {
				ActiveNode = tree[ActiveNode].slink
			} else {
				ActiveNode = Root
			}
		}
	}
}

func st_Show() {
	// using bfs to traverse tree

	//queue := make([]int, 0)
	//queue = append(queue, 1)
	//x = queue[0]
	//queue = queue[1:]

	nowIn := Root
	depth := 0

	for true {
		for _, child := range tree[nowIn].next { // process one node
			if child == Root {
				continue
			} // root can not be child

		}
	}
}

func main() {
	s := []byte("xabxa#babxba$")
	st_init()
	for _, char := range s {
		st_extend(char)
	}

}

type Node struct {
	start int
	end   int
	slink int
	next  [ALPHABET_SIZE]int
}
