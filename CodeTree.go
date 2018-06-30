package codetree

import (
	"fmt"
	"strings"
	"sync"
)

// Pool for CodeTree objects
var pool = sync.Pool{
	New: func() interface{} {
		return &CodeTree{}
	},
}

// CodeTree ...
type CodeTree struct {
	Line     string
	Children []*CodeTree
	Parent   *CodeTree
	Indent   int
}

// Close sends the tree and all of its children back to the memory pool.
// The resources are therefore freed up and the tree object should not be
// used after the final Close() call anymore.
func (tree *CodeTree) Close() {
	for _, child := range tree.Children {
		child.Close()
	}

	tree.Children = nil
	pool.Put(tree)
}

// New returns a tree structure if you feed it with indentantion based source code.
func New(src string) (*CodeTree, error) {
	ast := pool.Get().(*CodeTree)
	ast.Indent = -1
	ast.Line = ""
	ast.Parent = nil

	block := ast
	lastNode := ast
	lineStart := 0
	lineNumber := 0
	src = strings.Replace(src, "\r\n", "\n", -1)

	for i := 0; i <= len(src); i++ {
		if i != len(src) && src[i] != '\n' {
			continue
		}

		lineNumber++

		line := src[lineStart:i]
		lineStart = i + 1

		// Ignore empty lines
		empty := true

		for h := 0; h < len(line); h++ {
			if line[h] != '\t' && line[h] != ' ' {
				empty = false
				break
			}
		}

		if empty {
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

		if indent == block.Indent+1 {
			// OK
		} else if indent == block.Indent+2 {
			block = lastNode
		} else if indent <= block.Indent {
			for {
				block = block.Parent

				if block.Indent == indent-1 {
					break
				}
			}
		} else if indent > block.Indent+2 {
			return nil, fmt.Errorf("Invalid indentation at line: %s (%d)", line, lineNumber)
		}

		node := pool.Get().(*CodeTree)
		node.Line = line
		node.Indent = indent
		node.Parent = block
		lastNode = node
		block.Children = append(block.Children, node)
	}

	return ast, nil
}
