// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
)

// 代码生成器概要信息
type Profile struct {
	AppName string   // 应用名称
	Version string   // 应用版本
	Host    string   // 主机
	Port    int      // 端口，具化默认值
	User    string   // 用户，具化默认值
	Auth    string   // 密码，具化默认值
	DbName  string   // 数据库
	Ignores []string // 忽略的数据表
	Package string   // 包名
	Plump   bool     // 生成增删改查等数据操作代码
	Out     string   // 输出路径
	Pile    bool     // 单文件输出
}

func (c *Profile) String() string {
	return fmt.Sprintf("%+v", *c)
}
