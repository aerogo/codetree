package codetree_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aerogo/codetree"
)

func TestCodeTree(t *testing.T) {
	by, _ := ioutil.ReadFile("test/example.txt")
	code := string(by)
	tree, err := codetree.New(code)

	assert.NoError(t, err)
	defer tree.Close()

	assert.Equal(t, -1, tree.Indent)
	assert.Equal(t, 6, len(tree.Children))
	assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
	assert.Equal(t, "日本語", tree.Children[1].Children[0].Line)
}

func TestBadIndentation(t *testing.T) {
	by, _ := ioutil.ReadFile("test/bad-indentation.txt")
	code := string(by)
	tree, err := codetree.New(code)

	assert.Nil(t, tree)
	assert.Error(t, err)
}

func BenchmarkCodeTree(b *testing.B) {
	by, _ := ioutil.ReadFile("test/example.txt")
	code := string(by)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tree, _ := codetree.New(code)
		tree.Close()
	}
}

func TestComments(t *testing.T) {
	reader, _ := os.Open("test/comments.txt")
	tree, err := codetree.FromReader(reader)
	defer reader.Close()
	assert.Nil(t, err)
	assert.Equal(t, codetree.RootType, tree.Type)
	assert.Equal(t, 9, len(tree.Children))
	assert.Equal(t, codetree.CommentType, tree.Children[6].Type)
	assert.Equal(t, codetree.LineType, tree.Children[7].Type)
	assert.Equal(t, codetree.CommentType, tree.Children[7].Children[0].Type)
	assert.Equal(t, codetree.LineType, tree.Children[7].Children[2].Type)
	assert.Equal(t, codetree.CommentType, tree.Children[7].Children[2].Children[0].Type)
	assert.Equal(t, codetree.LineType, tree.Children[7].Children[3].Type)
	assert.Equal(t, codetree.CommentType, tree.Children[7].Children[3].Children[0].Type)
	assert.Equal(t, codetree.LineType, tree.Children[7].Children[4].Type)
	assert.Equal(t, codetree.CommentType, tree.Children[7].Children[4].Children[0].Type)
	assert.Equal(t, 5, len(tree.Children[7].Children))
	assert.Equal(t, codetree.CommentType, tree.Children[8].Type)
}

func TestFromFilelist(t *testing.T) {
	list := []string{"test/example.txt", "test/comments.txt"}
	tree, err := codetree.FromFilelist(list)
	assert.Nil(t, err)
	assert.Equal(t, codetree.RootType, tree.Type)
	for i, filename := range list {
		assert.Equal(t, filename, tree.Children[i].Filename)
		assert.Equal(t, codetree.RootType, tree.Children[i].Type)
		assert.Equal(t, 23, tree.Children[i].Children[4].LineNumber)
		assert.Equal(t, 25, tree.Children[i].Children[4].Children[1].LineNumber)
		assert.Equal(t, 28, tree.Children[i].Children[4].Children[2].Children[1].LineNumber)
	}
	assert.Equal(t, 34, tree.Children[1].Children[6].LineNumber)
	assert.Equal(t, 35, tree.Children[1].Children[7].LineNumber)
	assert.Equal(t, 36, tree.Children[1].Children[7].Children[0].LineNumber)
	assert.Equal(t, 37, tree.Children[1].Children[7].Children[1].LineNumber)

	assert.Equal(t, 38, tree.Children[1].Children[7].Children[2].LineNumber)
	assert.Equal(t, 38, tree.Children[1].Children[7].Children[2].Children[0].LineNumber)

	assert.Equal(t, 39, tree.Children[1].Children[7].Children[3].LineNumber)
	assert.Equal(t, 39, tree.Children[1].Children[7].Children[3].Children[0].LineNumber)

	assert.Equal(t, 40, tree.Children[1].Children[7].Children[4].LineNumber)
	assert.Equal(t, 40, tree.Children[1].Children[7].Children[4].Children[0].LineNumber)

	assert.Equal(t, 47, tree.Children[1].Children[8].LineNumber)
}
