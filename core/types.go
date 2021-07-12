// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package core

// type Type string
//
// const (
// 	Bool        Type = "bool"
// 	Byte        Type = "byte"
// 	ByteArray   Type = "[]byte"
// 	Int8        Type = "int8"
// 	Int16       Type = "int16"
// 	Int32       Type = "int32"
// 	Int64       Type = "int64"
// 	Float32     Type = "Float32"
// 	Float64     Type = "Float64"
// 	Money       Type = "int64"
// 	Decimal     Type = "decimal.Decimal" // "github.com/shopspring/decimal"
// 	DateTime    Type = "time.Time"
// 	Int64Array  Type = "pq.Int64Array"
// 	StringArray Type = "pq.StringArray"
// 	String      Type = "string"
// 	Interface   Type = "interface{}"
// )


//
// //类型转换，没有的类型在这里面添加
// func typeConvert(s string) string {
// 	if strings.Contains(s, "[]") {
// 		if strings.Contains(s, "char") || strings.Contains(s, "text") {
// 			return "pq.StringArray"
// 		}
// 		if strings.Contains(s, "integer") {
// 			return "pq.Int64Array"
// 		}
// 	}
// 	if strings.Contains(s, "char") || In(s, []string{"text", "longtext"}) {
// 		return "string"
// 	}
// 	if In(s, []string{"bigserial", "serial", "big serial", "int"}) {
// 		return "int"
// 	}
// 	if In(s, []string{"bigint"}) {
// 		return "int64"
// 	}
// 	if In(s, []string{"integer"}) {
// 		return "int32"
// 	}
// 	if In(s, []string{"smallint"}) {
// 		return "int16"
// 	}
// 	if In(s, []string{"numeric", "decimal", "real"}) {
// 		return "decimal.Decimal"
// 	}
// 	if In(s, []string{"bytea"}) {
// 		return "[]byte"
// 	}
// 	if strings.Contains(s, "time") || In(s, []string{"date"}) {
// 		return "time.Time"
// 	}
//
// 	return "interface{}"
// }