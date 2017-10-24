package codetree_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aerogo/codetree"
)

func TestCodeTree(t *testing.T) {
	bytes, _ := ioutil.ReadFile("example.txt")
	code := string(bytes)
	tree, err := codetree.New(code)

	assert.NoError(t, err)
	assert.Equal(t, -1, tree.Indent)
	assert.Equal(t, 6, len(tree.Children))
	assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
}

func BenchmarkCodeTree(b *testing.B) {
	bytes, _ := ioutil.ReadFile("example.txt")
	code := string(bytes)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		codetree.New(code)
	}
}
