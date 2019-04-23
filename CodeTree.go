package codetree

import (
	"fmt"
	"io"
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
	Line       string
	Children   []*CodeTree
	Parent     *CodeTree
	Indent     int
	LineNumber int
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

// New returns a tree structure if you feed it with indentation based source code.
func New(src string) (*CodeTree, error) {
	ast := pool.Get().(*CodeTree)
	ast.Indent = -1
	ast.Line = ""
	ast.Parent = nil
	ast.LineNumber = 0

	block := ast
	lastNode := ast
	lineStart := 0
	numLines := 0
	lineNumber := 1
	src = strings.Replace(src, "\r\n", "\n", -1)
	srcLen := len(src)
	comSrcLen := srcLen - 2

	for i := 0; i <= srcLen; i++ {

		line := ""
		lineNumber += numLines

		// Find multiline comments
		comOpenEnd := i + 2
		if i < comSrcLen && src[i:comOpenEnd] == "/*" {
			comLen := strings.Index(src[comOpenEnd:], "*/")
			if comLen == -1 {
				i = srcLen - 1
			} else {
				i += comLen + 3
			}
			line = src[lineStart : i+1]
			numLines = strings.Count(line, "\n")

		} else if i != srcLen && src[i] != '\n' { // or skip forward until a full line is found
			numLines = 0
			continue

		} else { // line found
			line = src[lineStart:i]
			numLines = 1
		}

		lineStart = i + 1

		// Ignore empty lines
		empty := true
		tabs := 0
		spaces := 0
		h := 0

	loop:
		for ; h < len(line); h++ {
			switch line[h] {
			case '\t':
				tabs++
			case ' ':
				spaces++
			default:
				empty = false
				break loop
			}
		}

		if empty {
			continue
		}

		// Indentation
		indent := tabs + spaces/2

		if h != 0 {
			line = line[h:]
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
			return nil, fmt.Errorf("Invalid indentation on line %d: %s", lineNumber, line)
		}

		node := pool.Get().(*CodeTree)
		node.Line = line
		node.Indent = indent
		node.Parent = block
		node.LineNumber = lineNumber
		lastNode = node
		block.Children = append(block.Children, node)
	}

	return ast, nil
}

const (
	lineStartState = iota
	lineState
	commentStartState
	commentLineState
	multilineCommentState
	multilineCommentEndState
)

func NewFromReader(src io.Reader) (*CodeTree, error) {
	ast := pool.Get().(*CodeTree)
	ast.Indent = -1
	ast.Line = ""
	ast.Parent = nil
	ast.LineNumber = 0

	block := ast
	lastNode := ast

	var (
		lineNumber int
		b          byte
		tabs       int
		spaces     int
		state      int
		err        error
	)

	readBytes := 512

	line := make([]byte, 0, 128)
	buf := make([]byte, readBytes)

	endLine := func() error {
		state = lineStartState

		// Ignore empty lines
		if len(line) == 0 {
			return nil
		}

		// Indentation
		indent := tabs + spaces/2

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
			return fmt.Errorf("invalid indentation on line %d: %s", lineNumber, line)
		}

		node := pool.Get().(*CodeTree)
		node.Line = string(line)
		node.Indent = indent
		node.Parent = block
		node.LineNumber = lineNumber
		lastNode = node
		block.Children = append(block.Children, node)

		// Reset
		tabs = 0
		spaces = 0
		line = line[:0]

		return nil
	}

	for {
		br, _ := src.Read(buf)
		for i := 0; i < br; i++ {
			b = buf[i]

			// Match states first, then tokens.
			// More verbose but, with fewer possible states than possible tokens, optimal.
			//*
			if b == '\r' {
				continue
			}

			switch state {
			case lineStartState:
				switch b {
				default:
					state = lineState
					line = append(line, b)
				case '\t':
					tabs++
				case ' ':
					spaces++
				case '\n':
					lineNumber++
					if err = endLine(); err != nil {
						return nil, err
					}
				case '/':
					if err = endLine(); err != nil {
						return nil, err
					}
					state = commentStartState
					line = append(line, b)
				}
			case lineState:
				switch b {
				default:
					line = append(line, b)
				case '\n':
					lineNumber++
					if err = endLine(); err != nil {
						return nil, err
					}
				case '/':
					if err = endLine(); err != nil {
						return nil, err
					}
					state = commentStartState
					line = append(line, b)
				}
			case commentStartState:
				switch b {
				default:
					state = lineState
				case '/':
					state = commentLineState
				case '*':
					state = multilineCommentState
				case '\n':
					lineNumber++
					if err = endLine(); err != nil {
						return nil, err
					}
				}
				line = append(line, b)
			case commentLineState:
				switch b {
				default:
					line = append(line, b)
				case '\n':
					lineNumber++
					if err = endLine(); err != nil {
						return nil, err
					}
				}
			case multilineCommentState:
				switch b {
				case '*':
					state = multilineCommentEndState
				case '\n':
					lineNumber++
				}
				line = append(line, b)
			case multilineCommentEndState:
				switch b {
				case '\n':
					lineNumber++
					fallthrough
				default:
					state = multilineCommentState
					line = append(line, b)
				case '/':
					line = append(line, b)
					if err = endLine(); err != nil {
						return nil, err
					}
				}
			}
			/*/
			switch b {
			case '\r':
				continue
			case '\t':
				switch state {
				case lineStartState:
					tabs++
				default:
					line = append(line, b)
				}
			case ' ':
				switch state {
				case lineStartState:
					spaces++
				default:
					line = append(line, b)
				}
			default:
				switch state {
				case lineStartState, commentStartState:
					state = lineState
				case multilineCommentEndState:
					state = multilineCommentState
				}
				line = append(line, b)
			case '\n':
				lineNumber++
				switch state {
				case lineStartState, lineState, commentLineState, commentStartState:
					if err := endLine(); err != nil {
						return nil, err
					}
				case multilineCommentState:
					line = append(line, b)
				case multilineCommentEndState:
					state = multilineCommentState
					line = append(line, b)
				}
			case '/':
				switch state {
				case lineStartState, lineState:
					if err := endLine(); err != nil {
						return nil, err
					}
					state = commentStartState
				case commentStartState:
					state = commentLineState
				case multilineCommentEndState:
					line = append(line, b)
					if err := endLine(); err != nil {
						return nil, err
					}
					continue
				}
				line = append(line, b)
			case '*':
				switch state {
				case commentStartState:
					state = multilineCommentState
				case multilineCommentState:
					state = multilineCommentEndState
				}
				line = append(line, b)
			}
			//*/
		}
		if br < readBytes {
			break
		}
	}
	lineNumber++
	if err = endLine(); err != nil {
		return nil, err
	}
	return ast, nil
}
