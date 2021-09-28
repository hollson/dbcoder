package internal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hollson/dbcoder/dialect/pgsql"
	schema2 "github.com/hollson/dbcoder/schema"
	"github.com/hollson/dbcoder/utils"
)

const (
	AppName = "dbcoder"
	VERSION = "v1.0.0"
)

type Generator struct {
	schema2.Driver // 数据库驱动类型
	Profile        // 配置文件
}

// Determine whether the driver is supported
func (g *Generator) Supported() bool {
	ds := []schema2.Driver{
		schema2.MySQL,
		schema2.PostgreSQL,
		// SQLite,
		// MongoDB,
		// MariaDB,
		// Oracle,
		// SQLServer,
	}
	for _, v := range ds {
		//		fmt.Println("=======================> ", v, g.Driver)
		if v == g.Driver {
			return true
		}
	}
	return false
}

// 执行生成命令
func (g *Generator) Generate() error {
	schema := g.factory()
	tables, err := schema.Tables()
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return errors.New("the count of tables in the database is 0")
	}

	if err := os.MkdirAll(g.Out, os.ModePerm); err != nil {
		return err
	}

	// 单文件输出
	if g.Pile {
		gofile := utils.PathTrim(fmt.Sprintf("%s/%s.go", g.Out, g.DbName))
		f, err := os.Create(gofile)
		if err != nil {
			return err
		}
		data := g.template(schema, tables...)
		if err := Execute(f, data); err != nil {
			return err
		}
		fmt.Printf(" 📖 生成文件：%s\n", gofile)
		return nil
	}

	// 多文件输出
	for _, table := range tables {
		gofile := utils.PathTrim(fmt.Sprintf("%s/%s.go", g.Out, table.Name))
		f, err := os.Create(gofile)
		if err != nil {
			return err
		}

		data := g.template(schema, table)
		if err := Execute(f, data); err != nil {
			return err
		}
		fmt.Printf(" 🚀 生成文件: %s\n", gofile)
	}
	return nil
}

// 生成器工厂 gen
func (g *Generator) factory() schema2.Schema {
	switch g.Driver {
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
		return pgsql.New(g.Host, g.Port, g.User, g.Auth, g.DbName, g.Ignores)
	}
}

func (g *Generator) template(schema schema2.Schema, tables ...schema2.Table) *GenTemplate {
	t := &GenTemplate{
		Generator: g.AppName,
		Version:   g.Version,
		Source:    fmt.Sprintf("%s://%s:%d/%s", g.Driver, g.Host, g.Port, g.DbName),
		Date:      time.Now().Format("2006-01-02"),
		Package:   g.Package,
	}
	if len(tables) == 1 {
		t.Source = fmt.Sprintf("%s://%s:%d/%s/%s", g.Driver, g.Host, g.Port, g.DbName, tables[0].Name)
	}

	for _, table := range tables {
		obj := Struct{
			Name:    utils.Pascal(table.Name),
			Comment: table.Comment,
		}
		for _, column := range table.Columns {
			obj.Fields = append(obj.Fields, Field{
				Name:    utils.Pascal(column.Name),
				Type:    schema.TypeMapping(column.Type).Name,
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
