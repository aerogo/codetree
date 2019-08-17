package codetree_test

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/aerogo/codetree"
	"github.com/akyoto/assert"
)

func TestCodeTree(t *testing.T) {
	// Repeat this test multiple times to detect memory corruption errors.
	// Incorrect usage of the "unsafe" package can lead to those errors.
	for i := 0; i < 1000; i++ {
		file, err := os.Open("testdata/example.txt")
		assert.Nil(t, err)
		defer file.Close()

		tree, err := codetree.New(file)
		assert.Nil(t, err)
		defer tree.Close()

		assert.Equal(t, -1, tree.Indent)
		assert.Equal(t, 6, len(tree.Children))
		assert.Equal(t, 1, len(tree.Children[5].Children))
		assert.Equal(t, "parent1", tree.Children[0].Line)
		assert.Equal(t, "parent2", tree.Children[1].Line)
		assert.Equal(t, "parent3", tree.Children[2].Line)
		assert.Equal(t, "child1", tree.Children[0].Children[0].Line)
		assert.Equal(t, "child1", tree.Children[1].Children[0].Line)
		assert.Equal(t, "child1", tree.Children[2].Children[0].Line)
		assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
	}

}

func TestWindowsLineEndings(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/example.txt")
	assert.Nil(t, err)
	code := string(data)
	code = strings.ReplaceAll(code, "\n", "\r\n")
	tree, err := codetree.New(strings.NewReader(code))
	assert.Nil(t, err)
	defer tree.Close()

	assert.Equal(t, -1, tree.Indent)
	assert.Equal(t, 6, len(tree.Children))
	assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
}

func TestBadIndentation(t *testing.T) {
	file, err := os.Open("testdata/bad-indentation.txt")
	assert.Nil(t, err)
	defer file.Close()

	tree, err := codetree.New(file)
	assert.Nil(t, tree)
	assert.NotNil(t, err)
}

func TestLongExample(t *testing.T) {
	file, err := os.Open("testdata/long-example.txt")
	assert.Nil(t, err)
	defer file.Close()

	tree, err := codetree.New(file)
	assert.Nil(t, err)
	defer tree.Close()

	assert.Equal(t, -1, tree.Indent)
	assert.Equal(t, "parent1", tree.Children[0].Line)
	assert.Equal(t, "parent2", tree.Children[1].Line)
	assert.Equal(t, "parent3", tree.Children[2].Line)
	assert.Equal(t, "child1", tree.Children[0].Children[0].Line)
	assert.Equal(t, "child1", tree.Children[1].Children[0].Line)
	assert.Equal(t, "child1", tree.Children[2].Children[0].Line)
	assert.Equal(t, "child1", tree.Children[5].Children[0].Line)
}

func TestTimeoutReader(t *testing.T) {
	file, err := os.Open("testdata/example.txt")
	assert.Nil(t, err)
	defer file.Close()
	timeoutReader := iotest.TimeoutReader(file)

	tree, err := codetree.New(timeoutReader)
	assert.NotNil(t, err)
	assert.Nil(t, tree)
}

func BenchmarkCodeTree(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	file, err := os.Open("testdata/example.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for i := 0; i < b.N; i++ {
		_, err = file.Seek(0, io.SeekStart)

		if err != nil {
			panic(err)
		}

		tree, err := codetree.New(file)

		if err != nil {
			panic(err)
		}

		tree.Close()
	}
}
