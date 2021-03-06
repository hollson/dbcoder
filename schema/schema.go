// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schema

// 定义了代码生成器所依赖的数据库信息
type Schema interface {
	Tables() ([]Table, error)
	TypeMapping(string) Type
	// 类型映射
}

type SchemaType string

// Go结构体字段类型
type Type struct {
	// 类型名称
	//  如：int，string，time.Time，decimal.Decimal等
	Name string
	// 引用的包
	//  如：time，github.com/shopspring/decimal等
	Pack string
}

// 数据库表
type Table struct {
	Name        string   // 表名
	Columns     []Column // 字段
	Comment     string   // 注释
	SpecialPack []string // 特殊类型的依赖包，如pg.Int64Array类型,须引用"github.com/lib/pg"
}

// 数据库表字段
type Column struct {
	Primary bool   // 是否主键
	Name    string // 字段
	Type    string // 类型
	Size    int    // 长度
	NotNull bool   // 不为Null
	Default string // 默认值
	Comment string // 注释
}
