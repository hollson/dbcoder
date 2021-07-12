// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pgsql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hollson/dbcoder/core"
	"github.com/hollson/dbcoder/utils"
)

const Empty = ""

type Database struct {
	Host   string
	Port   int
	User   string
	Auth   string
	DbName string
}

func New(c *core.Config) *Database {
	gen := &Database{
		Host:   c.Host,
		Port:   c.Port,
		User:   c.User,
		Auth:   c.Auth,
		DbName: c.DbName,
	}
	if gen.Port == 0 {
		gen.Port = 5432
	}
	if len(gen.User) == 0 {
		gen.User = "postgres"
	}
	if len(gen.Auth) == 0 {
		gen.Auth = "postgres"
	}
	return gen
}

// 连接字符串
func (g *Database) source() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		g.User, g.Auth, g.Host, g.Port, g.DbName)
}

// 查询数据库表清单SQL
func (g *Database) tablesSQL() string {
	return `SELECT a.tablename,
			COALESCE(c.description,'') AS comment
			FROM pg_tables a
			LEFT JOIN pg_class b on a.tablename=b.relname
			LEFT JOIN pg_description c on  b.oid=c.objoid and c.objsubid=0
			WHERE a.schemaname='public';`
}

// 查询数据表定义SQL
func (g *Database) columnsSQL(tableName string) string {
	return fmt.Sprintf(`
SELECT  a.attname AS field_name,	--字段表名
		a.attnotnull AS not_null,	--是否为NULL
		a.attlen AS field_size,		-- 字段大小
		COALESCE (ct.contype = 'p', FALSE ) AS is_primary_key,				-- 是否主键
		COALESCE (pg_get_expr(ad.adbin, ad.adrelid),'') AS default_value,	-- 默认值
		COALESCE(b.description,'') AS comment,								--注释
		CASE WHEN a.atttypid = ANY ('{int,int8,int2}'::regtype[]) AND EXISTS (SELECT 1 FROM pg_attrdef ad WHERE ad.adrelid = a.attrelid AND ad.adnum = a.attnum )
			THEN CASE a.atttypid
				WHEN 'int'::regtype THEN 'serial'
				WHEN 'int8'::regtype THEN 'bigserial'
				WHEN 'int2'::regtype THEN 'smallserial' END
			WHEN a.atttypid = ANY ('{uuid}'::regtype[]) AND COALESCE (pg_get_expr(ad.adbin, ad.adrelid ),'')<>''
				THEN 'autogenuuid' ELSE format_type( a.atttypid, a.atttypmod )
		END AS field_type										-- 标识类型 
FROM pg_attribute a
	INNER JOIN ONLY pg_class C ON C.oid = a.attrelid
	INNER JOIN ONLY pg_namespace n ON n.oid = C.relnamespace
	LEFT JOIN pg_constraint ct ON ct.conrelid = C.oid AND a.attnum = ANY ( ct.conkey ) AND ct.contype = 'p'
	LEFT JOIN pg_attrdef ad ON ad.adrelid = C.oid AND ad.adnum = a.attnum
	LEFT JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid
	LEFT join pg_type t ON a.atttypid = t.oid
WHERE a.attisdropped = FALSE AND a.attnum > 0 AND n.nspname = 'public' AND C.relname ='%s' -- 表名
ORDER BY a.attnum
`, tableName)
}

//
func (g *Database) Tables() (ret []core.Table, err error) {
	_db, err := sql.Open("postgres", g.source())
	if err != nil {
		return nil, err
	}
	fmt.Printf(" 💻 连接数据库: %s\n", g.source())
	defer _db.Close()

	rows, err := _db.Query(g.tablesSQL())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t = core.Table{}
		if err := rows.Scan(&t.Name, &t.Comment); err != nil {
			return nil, err
		}
		cs, err := g.columns(t.Name, _db)
		if err != nil {
			return nil, err
		}
		t.Columns = cs
		ret = append(ret, t)
	}
	return
}

func (g *Database) columns(tableName string, db *sql.DB) (ret []core.Column, err error) {
	rows, err := db.Query(g.columnsSQL(tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t = core.Column{}
	for rows.Next() {
		if err := rows.Scan(
			&t.Name,
			&t.NotNull,
			&t.Size,
			&t.Primary,
			&t.Default,
			&t.Comment,
			&t.Type); err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}
	return
}

// 类型映射
func (g *Database) TypeMapping(_type string) core.Type {
	switch {
	case utils.ContainString(_type, "bigserial", "serial", "big serial", "int"):
		return core.Type{Name: "int", Pack: Empty}
	case _type == "smallint":
		return core.Type{Name: "int16", Pack: Empty}
	case _type == "integer":
		return core.Type{Name: "int32", Pack: Empty}
	case _type == "bigint":
		return core.Type{Name: "int64", Pack: Empty}
	case utils.ContainString(_type, "numeric", "decimal", "real"):
		return core.Type{Name: "decimal.Decimal", Pack: "github.com/shopspring/decimal"}
	case strings.Contains(_type, "time") || _type == "date":
		return core.Type{Name: "time.Time", Pack: "time"}
	case _type == "bytea":
		return core.Type{Name: "[]byte", Pack: Empty}
	case strings.Contains(_type, "char") || utils.ContainString(_type, "text", "longtext"):
		return core.Type{Name: "string", Pack: Empty}
	case strings.Contains(_type, "char") || strings.Contains(_type, "text"):
		return core.Type{Name: "pq.StringArray", Pack: "pq.StringArray"}
	case strings.Contains(_type, "integer"):
		return core.Type{Name: "pq.Int64Array", Pack: "pq.StringArray"}
	default:
		return core.Type{Name: "interface{}", Pack: Empty}
	}
}
