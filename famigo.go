package main

import (
	"github.com/famigo/lang/compiler"
	"fmt"
)

func main() {
	compiler.CompilePackage("github.com/famigo/example/hello")
	fmt.Println("done")
}
