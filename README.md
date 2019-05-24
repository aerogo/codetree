# codetree

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Sponsor][sponsor-image]][sponsor-url]

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

## Style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Sponsors

| [![Cedric Fung](https://avatars3.githubusercontent.com/u/2269238?s=70&v=4)](https://github.com/cedricfung) | [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) | [![Eduard Urbach](https://avatars3.githubusercontent.com/u/438936?s=70&v=4)](https://twitter.com/eduardurbach) |
| --- | --- | --- |
| [Cedric Fung](https://github.com/cedricfung) | [Scott Rayapoullé](https://github.com/soulcramer) | [Eduard Urbach](https://eduardurbach.com) |

Want to see [your own name here?](https://github.com/users/akyoto/sponsorship)

[godoc-image]: https://godoc.org/github.com/aerogo/codetree?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/codetree
[report-image]: https://goreportcard.com/badge/github.com/aerogo/codetree
[report-url]: https://goreportcard.com/report/github.com/aerogo/codetree
[tests-image]: https://cloud.drone.io/api/badges/aerogo/codetree/status.svg
[tests-url]: https://cloud.drone.io/aerogo/codetree
[coverage-image]: https://codecov.io/gh/aerogo/codetree/graph/badge.svg
[coverage-url]: https://codecov.io/gh/aerogo/codetree
[sponsor-image]: https://img.shields.io/badge/github-donate-green.svg
[sponsor-url]: https://github.com/users/akyoto/sponsorship
