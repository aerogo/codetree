package codetree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// Pool for CodeTree objects
var pool = sync.Pool{
	New: func() interface{} {
		return &CodeTree{}
	},
}

const (
	RootType = iota
	LineType
	CommentType
)

// CodeTree ...
type CodeTree struct {
	Line       string
	Children   []*CodeTree
	Parent     *CodeTree
	Root       *CodeTree
	Indent     int
	LineNumber int
	Filename   string
	Type       int
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

// String returns a string describing this node which can be used in error messages.
func (tree *CodeTree) String() string {
	s := tree.GetFilename()
	if s != "" {
		s = fmt.Sprintf("in file %s ", s)
	}
	if tree.LineNumber != 0 {
		s = fmt.Sprintf("%son line %d ", s, tree.LineNumber)
	}
	switch tree.Type {
	case RootType:
		s = fmt.Sprintf("%s[root]", s)
	case LineType:
		s = fmt.Sprintf("%s[line]", s)
	case CommentType:
		s = fmt.Sprintf("%s[comment]", s)
	}
	if tree.Line != "" {
		s = fmt.Sprintf("%s: %s", s, tree.Line)
	}
	return s
}

func (tree *CodeTree) GetFilename() string {
	if tree.Filename != "" {
		return tree.Filename
	}
	if tree.Root != nil {
		return tree.Root.Filename
	}
	return ""
}

// New returns a tree structure if you feed it with indentation based source code.
func New(src string) (*CodeTree, error) {
	reader := strings.NewReader(src)
	return FromReader(reader)
}

const (
	lineStartState = iota
	lineState
	commentStartState
	commentLineState
	multilineCommentState
	multilineCommentEndState
)

// FromReader returns a CodeTree, taking an io.Reader as an argument instead of a string.
// This approach is more idiomatic, versatile and efficient.
// Uses a finite state machine type pattern to parse src in a single pass.
func FromReader(src io.Reader) (*CodeTree, error) {
	ast := pool.Get().(*CodeTree)
	ast.Indent = -1
	ast.Line = ""
	ast.Parent = nil
	ast.Root = ast
	ast.LineNumber = 0
	ast.Filename = ""
	ast.Type = RootType

	block := ast
	lastNode := ast

	var (
		b      byte
		tabs   int
		spaces int
		state  int
		err    error
	)

	readBytes := 512
	nodeType := LineType
	lineNumber := 1
	lastLineNumber := 1
	line := make([]byte, 0, 128)
	buf := make([]byte, readBytes)

	addNode := func() error {
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
		node.Root = ast
		node.LineNumber = lastLineNumber
		node.Type = nodeType
		lastNode = node
		block.Children = append(block.Children, node)

		// Reset
		tabs = 0
		spaces = 0
		line = line[:0]
		nodeType = LineType
		return nil
	}

	endLine := func() error {
		if err = addNode(); err != nil {
			return err
		}
		lastLineNumber = lineNumber
		return nil
	}

	beginComment := func() error {
		indent := 0
		if len(line) > 0 {
			indent = tabs + spaces/2 + 1
		}
		if err = addNode(); err != nil {
			return err
		}
		tabs += indent
		nodeType = CommentType
		return nil
	}

	for {
		br, _ := src.Read(buf)
		for i := 0; i < br; i++ {
			b = buf[i]

			if b == '\r' {
				continue
			}

			// Match states first, then tokens.
			// More verbose but, with fewer possible states than possible tokens, optimal.
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
					state = commentStartState
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
					state = commentStartState
				}
			case commentStartState:
				switch b {
				default:
					state = lineState
				case '/':
					if err = beginComment(); err != nil {
						return nil, err
					}
					state = commentLineState
				case '*':
					if err = beginComment(); err != nil {
						return nil, err
					}
					state = multilineCommentState
				}
				line = append(line, '/', b)
				if b == '\n' {
					lineNumber++
					if err = endLine(); err != nil {
						return nil, err
					}
				}
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

// FromFilelist returns a CodeTree whose every child element is a CodeTree built from the contents of a single file as
// enumerated in filelist. Each child CodeTree has its Filename element set to the corresponding element in filelist.
// This structure is parsed properly by Scarlett.
func FromFilelist(filelist []string) (*CodeTree, error) {
	ast := pool.Get().(*CodeTree)
	ast.Indent = -1
	ast.Line = ""
	ast.Parent = nil
	ast.Children = make([]*CodeTree, len(filelist))

	var (
		wg       sync.WaitGroup
		errsLock sync.Mutex
		errs     []string
	)

	recordError := func(message string) {
		errsLock.Lock()
		errs = append(errs, message)
		errsLock.Unlock()
		return
	}

	for i, filename := range filelist {
		wg.Add(1)
		go func(i int, filename string) {
			file, err := os.Open(filename)
			if err != nil {
				recordError(fmt.Sprintf("error reading file:%s", err))
				return
			}
			tree, err := FromReader(file)
			file.Close()
			if err != nil {
				recordError(err.Error())
				return
			}
			tree.Filename = filename
			ast.Children[i] = tree
			wg.Done()
		}(i, filename)
	}

	wg.Wait()
	if len(errs) != 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return ast, nil
}
