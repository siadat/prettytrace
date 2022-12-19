package main

import (
	"fmt"
	"strings"

	"github.com/siadat/prettytrace"
)

func main() {
	Example1()
	Example2()
}

func Example1() {
	prettytrace.Print()
}

func Example2() {
	defer func() {
		var p = recover()
		fmt.Println(p)
		prettytrace.Print()
	}()
	strings.Repeat(" ", -1)
}
