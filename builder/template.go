// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

/*
使用template模板方式，将Struct对象转换为go文件
*/
package builder

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/hollson/dbcoder/utils"
)

// Go文件模板
var _template = `// Code generated by {{.Generator}}. DO NOT EDIT.
// source: {{.Source}}
// {{.Generator}}: {{.Version}}
// date: {{.Date}}

package {{.Package}} {{ifImports .Imports}}
{{range .Structs}}{{ifComment .Comment}}
type {{.Name}} struct {
	{{range .Fields}}{{.Name}} {{.Type}}  $BACKQUOTE{{.Tag}}$BACKQUOTE {{if ne .Comment ""}} //{{.Comment}}{{end}}
{{end}}}
{{end}}
`

var GormTemplete = `func BeforunUpdaet()  {

}

func AfterUpdaet()  {

}`

type GenTemplate struct {
	Generator string   // 生成器名称
	Version   string   // 生成器版本
	Source    string   // 生成的来源(模板数据来源，如：192.168.0.10:5432/testdb/tableName)
	Date      string   // 生成日期
	Package   string   // 包名
	Imports   []string // 依赖包
	Structs   []Struct // 结构体
}

// Go文件中的结构体
type Struct struct {
	Name    string   // 结构体名称
	Comment string   // 结构体注释
	Fields  []Field  // 结构体字段
	Methods []Method // todo：结构体方法
}

// Go文件中的结构体字段
type Field struct {
	Name    string // 字段名称
	Type    string // 字段类型
	Tag     string // 字段标签
	Comment string // 字段注释
}

// todo：Go文件中的方法
type Method struct {
	// ...
}

// 自定义模板表达式：加载依赖包
func IfImports(imports []string) string {
	if len(imports) > 0 {
		utils.RangeStringsFunc(imports, func(s string) string {
			return fmt.Sprintf("\t\"%s\"", s)
		})
		return fmt.Sprintf(`
import (
%s
)`, strings.Join(imports, "\n"))
	}
	return ""
}

// 自定义模板表达式：成员注释
func IfComment(comment string) string {
	if len(comment) > 0 {
		return fmt.Sprintf("// %s", comment)
	}
	return ""
}

// 模板生成输出内容
func Execute(w io.Writer, tpl *GenTemplate) error {
	t := template.New("text")                                            // 定义模板对象
	t = t.Funcs(template.FuncMap{"ifImports": IfImports})                // 控制自定义元素
	t = t.Funcs(template.FuncMap{"ifComment": IfComment})                // 控制自定义元素
	t, err := t.Parse(strings.Replace(_template, "$BACKQUOTE", "`", -1)) // 处理反引号
	if err != nil {
		return err
	}

	if err := t.Execute(w, tpl); err != nil {
		return err
	}
	return nil
}
