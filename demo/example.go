package main

import "fmt"

func main() {
	fmt.Println("Text logger")
	exampleTextLogger()

	fmt.Println("\nContext logger")
	exampleContextLogger()

	fmt.Println("\nJSON logger")
	exampleJSONLogger()

	fmt.Println("\nHooks logger")
	exampleHooksLogger()
}
