package main

import (
	"github.com/famigo/lang/processor"
	// "fmt"
	// "go/build"
)


func main()  {
	// gopath := build.Default.GOPATH
	// filepath := fmt.Sprintf("%s/src/sandbox/someapp/someapp.go", gopath)
	// processor.ProcessFile(filepath)

	// processor.ProcessPackage("github.com/famigo/example/hello")
	
	processor.ProcessPackage("github.com/famigo/example/hello")
}