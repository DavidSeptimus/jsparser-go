package predicates

import sitter "github.com/smacker/go-tree-sitter"

type Node = sitter.Node

//NodeType returns Node predicate that checks whether a node's type matches any of the supplied types
func NodeType(ntypes ...string) func(*Node) bool {
	return func(n *Node) bool {
		for _, ntype := range ntypes {
			if ntype == n.Type() {
				return true
			}
		}
		return false
	}
}

//Chain returns a chain of *Node predicates where the first predicate is tested against the input *Node
//and subsequent predicates are evaluated against the node returned from nextFunc
func Chain(nextFunc func(*Node) *Node, predicates ...func(*Node) bool) (chainedPredicate func(*Node) bool) {
	chainedPredicate = func(startNode *Node) bool {
		currentNode := startNode
		for _, predicate := range predicates {
			if currentNode == nil {
				return false
			}
			if predicate(currentNode) == false {
				return false
			}
			currentNode = nextFunc(currentNode)
		}
		return true
	}
	return chainedPredicate
}

//IsInvocation checks if a node's next sibling in a Node of type "arguments"
func IsInvocation(n *Node) bool {
	if n.Parent() == nil {
		return false
	}
	psibling := n.Parent().NextSibling()
	return psibling != nil && psibling.Type() == "arguments"
}

//TextEquals returns a predicate that compares the supplied *Node's text representation
// to the value of the supplied text argument using the supplied Scanner
func TextEquals(text string, scanner Scanner) func(*Node) bool {
	return func(n *Node) bool {
		return text == scanner.Find(n)
	}
}

//Scanner is an interface representing the ability to find a Node's text representation
type Scanner interface {
	Find(*Node) string
}
