package todoparser

import "testing"

func TestTodo(t *testing.T) {

	// Go代码
	code := `
package main

import "fmt"

func main() {
	// TODO: 完成这个函数
	fmt.Println("Hello, world!")
}
`

	todo(code)

}
