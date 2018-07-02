# codetree

[![Godoc reference][godoc-image]][godoc-url]
[![Go report card][goreportcard-image]][goreportcard-url]
[![Tests][travis-image]][travis-url]
[![Code coverage][codecov-image]][codecov-url]
[![License][license-image]][license-url]

Parses indented code (Python, Pug, Stylus, Pixy, Scarlet, etc.) and returns a tree structure.

## Installation

```bash
go get github.com/aerogo/codetree
```

## Usage

```go
tree, err := codetree.New(source)
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

See [CodeTree](https://github.com/aerogo/codetree/blob/master/CodeTree.go#L17-L22) structure.

The root node always starts with `Indent` being `-1`.

## Author

| [![Eduard Urbach on Twitter](https://gravatar.com/avatar/16ed4d41a5f244d1b10de1b791657989?s=70)](https://twitter.com/eduardurbach "Follow @eduardurbach on Twitter") |
|---|
| [Eduard Urbach](https://eduardurbach.com) |

[godoc-image]: https://godoc.org/github.com/aerogo/codetree?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/codetree
[goreportcard-image]: https://goreportcard.com/badge/github.com/aerogo/codetree
[goreportcard-url]: https://goreportcard.com/report/github.com/aerogo/codetree
[travis-image]: https://travis-ci.org/aerogo/codetree.svg?branch=master
[travis-url]: https://travis-ci.org/aerogo/codetree
[codecov-image]: https://codecov.io/gh/aerogo/codetree/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/aerogo/codetree
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://github.com/aerogo/codetree/blob/master/LICENSE
