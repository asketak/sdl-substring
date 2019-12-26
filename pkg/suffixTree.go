package suffixTree

type SuffixTree struct {
	root Node
}

type Node struct {
	parent *Node
	childs []Node
	value  string
}
