package todoparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func todo(code string) {

	// 解析代码
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		fmt.Println("无法解析代码:", err)
		return
	}

	for _, v := range file.Comments {
		fmt.Println("++++++++++>>>>> ", v.Text(), v.Pos())
	}

	// 遍历抽象语法树
	ast.Inspect(file, func(node ast.Node) bool {
		// 检查注释节点
		comment, ok := node.(*ast.Comment)
		fmt.Println("--====>>>>>>>>>> ", ok, comment)
		if ok {
			// 提取TODO注释内容
			if strings.Contains(comment.Text, "TODO") {
				todo := strings.TrimPrefix(comment.Text, "// TODO:")
				todo = strings.TrimSpace(todo)
				fmt.Println("TODO:", todo)
			}
		}
		return true
	})
}
