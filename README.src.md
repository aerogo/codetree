# {name}

{go:header}

Parses indented code (Python, Pug, Stylus, Pixy, codetree, etc.) and returns a tree structure.

## Installation

```bash
go get github.com/aerogo/codetree
```

## Usage

```go
tree, err := codetree.New(reader)
defer tree.Close()
```

## Input

```
parent1
	child1
	child2
	child3
		child3.1
		child3.2
	child4

parent2
	child1
```

## Output

See [CodeTree](https://github.com/aerogo/codetree/blob/master/CodeTree.go#L23-L32) structure.

The root node always starts with `Indent` being `-1`.

{go:footer}
