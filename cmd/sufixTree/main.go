package sufixTree

var len_s int
var stringg string
var len_string int
var max_len int
var root InternalNode

type leafnode struct {
	From_first_node bool
}

func newLeafnode(from_first_node bool) *leafnode {
	return &leafnode{From_first_node: from_first_node}
}

func (l *leafnode) has_s_leaves() bool {
	return l.From_first_node
}

func (l *leafnode) has_t_leaves() bool {
	return !l.From_first_node
}

type InternalNode struct {
	edges map[string] Edge
	edgesKeys map[string]bool
	link *InternalNode
	root_lenght int
	has_s_leaves bool
	has_t_leaves bool
	already_counted bool
}

func NewInternalNode(root_len int) *InternalNode {
	return &InternalNode{
		link: nil,
		root_lenght:root_len,
	}
}

func (i *InternalNode) getitem(key string) Edge {
	return i.edges[key]
}

func (i *InternalNode) setitem(key string, edge Edge) {
	i.edges[key] = edge
	i.edgesKeys[key] = true
	i.has_s_leaves = i.has_s_leaves || edge.dest.has_s_leaves
	i.has_t_leaves = i.has_t_leaves || edge.dest.has_t_leaves
}

func (i *InternalNode) contains(key string) bool {
	return i.edgesKeys[key]
}

type Edge struct {
	dest   InternalNode
	start  int
	end    int
	length int
}

func NewEdge(dest InternalNode, start int, end int) *Edge {
	return &Edge{dest: dest, start: start, end: end, length: end - start	}
}

type Cursor struct {
	node InternalNode
	edge string
	idx int
	lag int
}

func NewCursor(node InternalNode) *Cursor {
	return &Cursor{node: node}
}

func (c *Cursor) is_followed_by(letter string) bool {
	if c.idx == 0 {
		return c.node.contains(letter)
	}
	x := c.node.getitem(c.edge).start + c.idx
	return letter == stringg[x:x+1]
}

func (c *Cursor) defr(letter string) {
	c.idx += 1

	if c.edge == "" {
		c.edge = letter
	}
	edge := c.node.getitem(c.edge)
	if c.idx == edge.length{
		c.node = edge.dest
		c.edge = ""
		c.idx = 0
	}
}

func (c *Cursor) post_insert(i int ) {
	c.lag -= 1
	if c.node.root_lenght == root.root_lenght { // TODO mozna blbost
		if c.idx > 1{
			c.edge = stringg[i - c.lag:i-c.lag+1]
			c.idx -= 1
		}else{
			c.idx = 0
			c.edge = ""
		}
	}

	if c.node.link == nil {
		c.node = root
	}else{
		c.node = *c.node.link
	}

	for ; len(c.edge)>0 && c.idx >= c.node.getitem(c.edge); {

	}

}



func LCSukonnen( s string, t string)  {
	len_s = len(s)
	stringg = s + '#' + t + '$'
	len_string = len(stringg)
	max_len = 0
	root = NewInternalNode(0)
}


