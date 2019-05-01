package codetree_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aerogo/codetree"
)

func TestCodeTree(t *testing.T) {
	files := []string{
		"testdata/example.txt",
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		assert.NoError(t, err)
		code := string(data)
		tree, err := codetree.New(code)
		assert.NoError(t, err)
		defer tree.Close()

		assert.Equal(t, -1, tree.Indent)
		assert.Equal(t, 6, len(tree.Children))
		assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
	}
}

func TestWindowsLineEndings(t *testing.T) {
	files := []string{
		"testdata/example.txt",
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		assert.NoError(t, err)
		code := string(data)
		code = strings.ReplaceAll(code, "\n", "\r\n")
		tree, err := codetree.New(code)
		assert.NoError(t, err)
		defer tree.Close()

		assert.Equal(t, -1, tree.Indent)
		assert.Equal(t, 6, len(tree.Children))
		assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
	}
}

func TestBadIndentation(t *testing.T) {
	bytes, _ := ioutil.ReadFile("testdata/bad-indentation.txt")
	code := string(bytes)
	tree, err := codetree.New(code)

	assert.Nil(t, tree)
	assert.Error(t, err)
}

func BenchmarkCodeTree(b *testing.B) {
	bytes, _ := ioutil.ReadFile("example.txt")
	code := string(bytes)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree, err := codetree.New(code)

		if err != nil {
			b.Fail()
		}

		tree.Close()
	}
}
