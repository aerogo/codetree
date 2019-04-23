package codetree_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aerogo/codetree"
)

func TestCodeTree(t *testing.T) {
	bytes, _ := ioutil.ReadFile("test/example.txt")
	code := string(bytes)
	tree, err := codetree.New(code)

	assert.NoError(t, err)
	defer tree.Close()

	assert.Equal(t, -1, tree.Indent)
	assert.Equal(t, 6, len(tree.Children))
	assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
}

func TestBadIndentation(t *testing.T) {
	bytes, _ := ioutil.ReadFile("test/bad-indentation.txt")
	code := string(bytes)
	tree, err := codetree.New(code)

	assert.Nil(t, tree)
	assert.Error(t, err)
}

func BenchmarkCodeTree(b *testing.B) {
	bytes, _ := ioutil.ReadFile("test/example.txt")
	code := string(bytes)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree, _ := codetree.New(code)
		tree.Close()
	}
}

func TestCodeTreeFromReader(t *testing.T) {
	file, _ := os.Open("test/example.txt")
	tree, err := codetree.NewFromReader(file)

	assert.NoError(t, err)
	defer tree.Close()

	assert.Equal(t, -1, tree.Indent)
	assert.Equal(t, 6, len(tree.Children))
	assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
}

func TestCodeTreeCompare(t *testing.T) {
	bytes, _ := ioutil.ReadFile("test/example.txt")
	code := string(bytes)
	tree, err := codetree.New(code)
	assert.NoError(t, err)
	defer tree.Close()

	file, _ := os.Open("test/example.txt")
	defer file.Close()
	rTree, err := codetree.NewFromReader(file)
	assert.NoError(t, err)
	defer rTree.Close()

	assert.Equal(t, tree, rTree)
}

func BenchmarkCodeTreeFromReader(b *testing.B) {
	by, _ := ioutil.ReadFile("test/example.txt")
	reader := bytes.NewReader(by)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		tree, _ := codetree.NewFromReader(reader)
		tree.Close()
	}
}

func BenchmarkCodeTreeOpen(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		by, _ := ioutil.ReadFile("test/example.txt")
		code := string(by)
		tree, _ := codetree.New(code)
		tree.Close()
	}
}

func BenchmarkCodeTreeFromReaderOpen(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		file, _ := os.Open("test/example.txt")
		tree, _ := codetree.NewFromReader(file)
		tree.Close()
		file.Close()
	}
}
