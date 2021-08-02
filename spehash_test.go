package spehash_test

import (
	"fmt"

	"github.com/anrid/spehash"
)

func ExampleHashtable_Insert() {
	t := spehash.NewTable()
	t.Insert("!!!!")
	t.Insert("!!!a")
	t.Insert("a")
	t.Insert("a")
	t.Insert("b")
	t.Insert("💩")

	fmt.Println(t.Dump())
	// Output: a101 b101
}

func ExampleHashtable_Insert_second() {
	t := spehash.NewTable()
	t.Insert("aaabbbcccddd")
	t.Insert("aaa")
	t.Insert("aaz")
	t.Insert("💩")

	fmt.Println(t.Dump())
	// Output: a101 z101
}

func ExampleHashtable_Delete() {
	t := spehash.NewTable()
	t.Insert("aaa")
	t.Insert("bbb")
	t.Insert("💩")

	// Bad deletes.
	t.Delete("!!!!")

	t.Delete("bbb")
	t.Delete("bbb")

	fmt.Println(t.Dump())
	// Output: a101 b011
}

func ExampleHashtable_Search() {
	t := spehash.NewTable()
	t.Insert("aaa")
	t.Insert("bbb")

	// Bad searches.
	t.Search("!!!!")
	t.Search("ccc")

	e, found := t.Search("bbb")

	fmt.Printf("%s %s %t\n", t.Dump(), e.Value, found)
	// Output: a101 b101 bbb true
}
