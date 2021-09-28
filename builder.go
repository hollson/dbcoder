// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// 从命令行或环境变量读取配置信息

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hollson/dbcoder/dialect/gorm"
	"github.com/hollson/dbcoder/internal"
	"github.com/hollson/dbcoder/internal/schema"
	"github.com/hollson/dbcoder/utils"
)

var (
	_driver  string
	_host    string
	_port    int
	_user    string
	_auth    string
	_dbName  string
	_package string
	_out     string
	_ignores string
	_version bool
	_help    bool
	_pile    bool
	_plump   bool
)

// 初始化Flag
func initFlag() {
	utils.Usage = Usage
	utils.StringVar(&_driver, "driver", "", "DB驱动")
	utils.StringVar(&_host, "host", "localhost", "主机名")
	utils.IntVar(&_port, "port", 0, "端口")
	utils.StringVar(&_user, "user", "", "用户名")
	utils.StringVar(&_auth, "auth", "", "密码")
	utils.StringVar(&_dbName, "dbname", "", "数据库名称")
	utils.BoolVar(&gorm.Gorm, "gorm", false, "是否添加gorm标签")
	utils.StringVar(&_out, "out", "./model", "输出路径")
	utils.StringVar(&_package, "package", "model", "go文件包名")
	utils.StringVar(&_ignores, "ignores", "", "忽略的表(用逗号分割,可使用通配符，如：t1,t2,t_*)")
	utils.BoolVar(&_plump, "plump", false, "生成CURD代码")
	utils.BoolVar(&_version, "version", false, "查看版本")
	utils.BoolVar(&_pile, "pile", false, "单文件输出")
	utils.BoolVar(&_help, "help", false, "查看帮助")
}

// 加载命令行参数
func Load() (*internal.Generator, error) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("\033[%dmAn error occurred: %v\033[0m\n\n", utils.FgRed, err)
	// 		Usage()
	// 		os.Exit(1)
	// 	}
	// }()

	initFlag()
	if err := utils.Parse(); err != nil {
		return nil, err
	}

	if _help || len(os.Args) == 1 {
		Usage()
		os.Exit(0)
	}

	if _version {
		fmt.Println(internal.VERSION)
		os.Exit(0)
	}

	if err := check(); err != nil {
		return nil, err
	}

	gen := &internal.Generator{
		Driver: schema.DriverValue(_driver),
		Profile: internal.Profile{
			AppName: internal.AppName,
			Version: internal.VERSION,
			Host:    _host,
			Port:    _port,
			User:    _user,
			Auth:    _auth,
			DbName:  _dbName,
			Package: _package,
			Ignores: strings.Split(_ignores, ","),
			Out:     _out,
			Pile:    _pile,
		},
	}
	if !gen.Supported() {
		return nil, fmt.Errorf("the driver named %v is not supported", _driver)
	}
	return gen, nil
}

// 检查命令行参数
func check() error {
	if len(_driver) == 0 {
		return fmt.Errorf("driver is needed")
	}

	if len(_dbName) == 0 {
		return fmt.Errorf("dbname is needed")
	}
	return nil
}

func Usage() {
	fmt.Println("\033[1;34m Welcome to dbcoder\033[0m")
	fmt.Printf("\u001B[1;35m       ____                   __         \n  ____/ / /_  _________  ____/ ___  _____\n / __  / __ \\/ ___/ __ \\/ __  / _ \\/ ___/\n/ /_/ / /_/ / /__/ /_/ / /_/ /  __/ /    \n\\__,_/_.___/\\___/\\____/\\__,_/\\___/_/     (%v)\u001B[0m\n", internal.VERSION)
	fmt.Printf(`
Usage:
    dbcoder <command> dbname=<dbName> [option]...

Command:
    mysql	从mysql数据库生成表实体
    pgsql	从postgres数据库生成表实体
    help	查看帮助

Option:
    -host	主机名
    -host	主机名
    -host	主机名
    -host	主机名
    -host	主机名
    -host	主机名

Default param:
    mysql: -host=localhost -port=3306 -user=root -auth=""
    pgsql: -host=localhost -port=5432 -user=postgres -auth=postgres

Example:
    dbcoder -driver=pgsql -dbname=testdb
    dbcoder -driver=pgsql -host=localhost -port=5432 -user=postgres -auth=postgres -dbname=testdb -gorm -package=entity

更多详情，请参考 https://github.com/hollson/dbcoder

`)
}
