package main

import (
	"context"
	"fmt"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"jsparser-go/predicates"
	"log"
	"os"
	"strings"
)

type Node = sitter.Node
type Tree = sitter.Tree

// prints any invocations of the property, propName, on the variable associated with the import of the module, moduleName,
// in the source file at srcPath
func main() {

	// inputs
	moduleName := "fs/promises"
	propName := "readFile"
	srcPath := "./resources/app.js"

	findInvocations(moduleName, propName, srcPath)
}

//findInvocations prints and returns the variable name associated with a module imported with require()
//as well as any invocations of the supplied method property name in the source file and their line numbers
func findInvocations(moduleName string, propName string, srcPath string) (mVar *Node, invocations []*Node) {
	src := readSource(srcPath)
	scanner := sourceScanner{&src}
	tree := parseJs(src)
	requiresExp := fmt.Sprintf("require(\"%s\")", moduleName)

	moduleAssignments := findModuleAssignment(requiresExp, tree, scanner)
	if len(moduleAssignments) == 0 {
		fmt.Printf("no import found for '%s'", moduleName)
		return
	}

	mVar = moduleAssignments[0] // just use the first assignment for now
	varName := scanner.Find(mVar)
	invocations = findPropReferences(varName, propName, tree, scanner)

	// just print here to keep things simple
	fmt.Printf("module \"%s\" imported as \"%s\"\n", moduleName, varName)
	fmt.Printf("all occurences of %s invoking %s:\n", varName, propName)
	printNodeLines(invocations, scanner)

	return
}

//findModuleAssignment returns a slice of type *Node containing all Nodes
//where the supplied module is assigned to a variable
func findModuleAssignment(module string, tree *Tree, scanner sourceScanner) []*Node {
	var results []*Node

	rootNode := tree.RootNode()
	nodes := findNode(rootNode, predicates.NodeType("variable_declarator", "assignment_expression"))
	for _, node := range nodes {
		callExps := findNode(node, predicates.NodeType("call_expression"))
		if len(callExps) == 0 {
			continue
		}
		callExp := callExps[0]
		if scanner.Find(callExp) == module {
			identifierNode := findNode(node, predicates.NodeType("identifier"))[0]
			results = append(results, identifierNode)
		}
	}
	return results
}

//printNodes prints each node in the supplied slice along with the line number of its start point
func printNodes(nodes []*Node, scanner sourceScanner) {
	for _, n := range nodes {
		fmt.Printf("line %d: %s\n", n.StartPoint().Row+1, scanner.Find(n))
	}
}

//printNodeLines prints the line number and line content associated with
//the start point of each node in the supplied slice
func printNodeLines(nodes []*Node, scanner sourceScanner) {
	lines := scanner.lines()
	for _, n := range nodes {
		fmt.Printf("line %d: %s\n", n.StartPoint().Row+1, lines[n.StartPoint().Row])
	}
}

//parseJs takes a javascript source file as a byte slice and returns tree-sitter Tree
func parseJs(src []byte) *Tree {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, src)
	if err != nil {
		log.Panicln(err)
	}

	return tree
}

//readSource returns byte slice containing the source file's content
func readSource(path string) []byte {
	src, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	return src
}

//children returns a slice containing reference to all the supplied Node's immediate children
func children(n *Node) []*Node {
	count := int(n.ChildCount())
	nodes := make([]*Node, count)

	for i := 0; i < count; i++ {
		nodes[i] = n.Child(i)
	}
	return nodes
}

//findPropReferences returns a slice of type *Node containing all references to property propName
// on variable varName
// note: current implementation has an additional restriction such that it only returns invocations rather than accesses
func findPropReferences(varName string, propName string, tree *Tree, scanner sourceScanner) []*Node {
	props := findNode(
		tree.RootNode(),
		predicates.Chain(prevSibling,
			predicates.Chain(identity,
				predicates.TextEquals(propName, scanner),
				predicates.NodeType("property_identifier"),
				predicates.IsInvocation, // this is probably an over-optimization),
			),
			predicates.NodeType("."),
			predicates.Chain(identity,
				predicates.TextEquals(varName, scanner),
				predicates.NodeType("identifier"),
			),
		),
	)

	return props
}

//identity returns the supplied *Node
func identity(n *Node) *Node {
	return n
}

//prevSibling returns the supplied *Node's previous sibling
func prevSibling(n *Node) *Node {
	return n.PrevSibling()
}

/*
findNode returns a slice of type *Node containing all Nodes that match the supplied predicate recursively;
starting from the root node (DFS - preorder traversal)

Note: might benefit from a negative predicate that actively excludes certain irrelevant branches
*/
func findNode(n *Node, predicate func(*Node) bool) []*Node {
	var results []*Node

	if predicate(n) == true {
		results = append(results, n)
	}

	for _, child := range children(n) {
		results = append(results, findNode(child, predicate)...)
	}

	return results
}

//sourceScanner provides methods for retrieving content from a source file's byte slice
type sourceScanner struct {
	Source *[]byte
}

//lines returns the source content as a slice of lines
func (s sourceScanner) lines() []string {
	return strings.Split(string(*s.Source), "\n")
}

//Find returns the text representation of a Node
func (s sourceScanner) Find(n *Node) string {
	return n.Content(*s.Source)
}
