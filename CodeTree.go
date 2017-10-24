package codetree

import (
	"errors"
	"strings"
)

// CodeTree ...
type CodeTree struct {
	Line     string
	Children []*CodeTree
	Parent   *CodeTree
	Indent   int
}

// New returns a tree structure if you feed it with indentantion based source code.
func New(src string) (*CodeTree, error) {
	ast := &CodeTree{
		Indent: -1,
	}

	block := ast
	lastNode := ast
	lineStart := 0
	src = strings.Replace(src, "\r\n", "\n", -1)

	for i := 0; i <= len(src); i++ {
		if i != len(src) && src[i] != '\n' {
			continue
		}

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

		node := &CodeTree{
			Line:   line,
			Indent: indent,
		}

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
			return nil, errors.New("Invalid indentation")
		}

		node.Parent = block
		block.Children = append(block.Children, node)

		lastNode = node
	}

	return ast, nil
}
