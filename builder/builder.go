// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package builder

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hollson/dbcoder/dialect/pgsql"
	"github.com/hollson/dbcoder/internal"
	"github.com/hollson/dbcoder/utils"
)

// dbcoder -d pgsql -h localhost -p 5432 -u postgres -auth 123456 -d deeplink -gorm

// 生成器工厂
func schemaFactory(driver internal.DatabaseDriver, c *internal.Config) internal.Schema {
	switch driver {
	// case core.MySQL:
	// 	return mysql.New(cfg)
	// case core.PostgreSQL:
	// 	return pgsql.New(cfg)
	// case core.SQLite:
	// 	return new(mysql.Generator)
	// case core.MariaDB:
	// 	return new(mysql.Generator)
	// case core.MongoDB:
	// 	return new(mysql.Generator)
	// case core.Oracle:
	// 	return new(mysql.Generator)
	default:
		return pgsql.New(c.Host, c.Port, c.User, c.Auth, c.DbName)
	}
}

// 执行生成命令
func Generate(driver internal.DatabaseDriver, cfg *internal.Config) error {
	schema := schemaFactory(driver, cfg)
	tables, err := schema.Tables()
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return errors.New("the count of tables in the database is 0")
	}

	if err := os.MkdirAll(cfg.Out, os.ModePerm); err != nil {
		return err
	}

	// 单文件输出
	if cfg.Pile {
		gofile := utils.PathTrim(fmt.Sprintf("%s/%s.go", cfg.Out, cfg.DbName))
		f, err := os.Create(gofile)
		if err != nil {
			return err
		}
		data := Schema2Template(driver, cfg, schema, tables...)
		if err := Execute(f, data); err != nil {
			return err
		}
		fmt.Printf(" 📖 生成文件：%s\n", gofile)
		return nil
	}

	// 多文件输出
	for _, table := range tables {
		gofile := utils.PathTrim(fmt.Sprintf("%s/%s.go", cfg.Out, table.Name))
		f, err := os.Create(gofile)
		if err != nil {
			return err
		}

		data := Schema2Template(driver, cfg, schema, table)
		if err := Execute(f, data); err != nil {
			return err
		}
		fmt.Printf(" 🚀 生成文件: %s\n", gofile)
	}
	return nil
}

func Schema2Template(driver internal.DatabaseDriver, cfg *internal.Config, schema internal.Schema, tables ...internal.Table) *GenTemplate {
	t := &GenTemplate{
		Generator: cfg.AppName,
		Version:   cfg.Version,
		Source:    fmt.Sprintf("%s://%s:%d/%s", driver, cfg.Host, cfg.Port, cfg.DbName),
		Date:      time.Now().Format("2006-01-02"),
		Package:   cfg.Package,
	}
	if len(tables) == 1 {
		t.Source = fmt.Sprintf("%s://%s:%d/%s/%s", driver, cfg.Host, cfg.Port, cfg.DbName, tables[0].Name)
	}

	for _, table := range tables {
		obj := Struct{
			Name:    utils.Pascal(table.Name),
			Comment: table.Comment,
		}
		for _, column := range table.Columns {
			obj.Fields = append(obj.Fields, Field{
				Name:    utils.Pascal(column.Name),
				Type:    schema.TypeMapping(column.Type).Name, // fixme 转换
				Tag:     column.Default,
				Comment: column.Comment,
			})
			if pack := schema.TypeMapping(column.Type).Pack; len(pack) > 0 {
				t.Imports = append(t.Imports, pack)
			}
		}
		t.Imports = utils.DistinctStringArray(t.Imports)
		t.Structs = append(t.Structs, obj)
	}
	return t
}
