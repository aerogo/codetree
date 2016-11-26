package codetree

import "strings"

// CodeTree ...
type CodeTree struct {
	Line     string
	Children []*CodeTree
	Parent   *CodeTree
	Indent   int
}

// New returns a tree structure if you feed it with indentantion based source code.
func New(src string) *CodeTree {
	ast := new(CodeTree)
	ast.Indent = -1

	block := ast
	lastNode := ast

	lines := strings.Split(strings.Replace(src, "\r\n", "\n", -1), "\n")

	for _, line := range lines {
		// Ignore empty lines
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Indentation
		indent := 0
		for indent < len(line) {
			if line[indent] != '\t' {
				break
			}

			indent++
		}

		if indent != 0 {
			line = line[indent:]
		}

		node := new(CodeTree)
		node.Line = line
		node.Indent = indent

		if node.Indent == block.Indent+1 {
			// OK
		} else if node.Indent == block.Indent+2 {
			block = lastNode
		} else if node.Indent <= block.Indent {
			for {
				block = block.Parent
				if block.Indent == node.Indent-1 {
					break
				}
			}

		} else if node.Indent > block.Indent+2 {
			panic("Invalid indentation")
		}

		node.Parent = block
		block.Children = append(block.Children, node)

		lastNode = node
	}

	return ast
}
