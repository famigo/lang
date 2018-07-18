package main

import (
	"fmt"
	"github.com/famigo/lang/processor"
)

func main() {
	processor.ProcessPackage("github.com/famigo/example/hello")
	fmt.Println("done")
}
