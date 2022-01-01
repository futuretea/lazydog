package main

import (
	"github.com/JodeZer/lazydog/example/expdir"
	"github.com/JodeZer/lazydog/example/expdir/expdir2"
	"github.com/JodeZer/lazydog/example/expdir/expdir2/expdir2_1"
	"github.com/JodeZer/lazydog/example/expdir/expdir2/expdir3"
	"github.com/JodeZer/lazydog/example/expdir/expdir2/expdir3/expdir4"
)

func main() {
	expdir.Foo()
	expdir2.Foo()
	expdir2_1.Foo()
	expdir3.Foo()
	expdir4.Foo()
}
