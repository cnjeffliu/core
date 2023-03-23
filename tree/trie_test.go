package tree_test

import (
	"fmt"
	"testing"

	"github.com/cnzf1/gocore/tree"
)

func TestTrie(t *testing.T) {
	tr := tree.NewTrie()
	tr.Insert("Hello")
	tr.Insert("Hello1")
	tr.Insert("Hello2")
	fmt.Println(tr.Search("Hello"))
	fmt.Println(tr.Search("Hallo"))
	fmt.Println(tr.Search("Hello2"))

	tr.Insert("河北")
	tr.Insert("湖南")
	tr.Insert("湖北")
	fmt.Println(tr.Search("湖北"))
	fmt.Println(tr.Search("湖北1"))
}
