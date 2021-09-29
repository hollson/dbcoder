// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package pgsql

import (
	"database/sql"
	"fmt"

	"github.com/hollson/dbcoder/schema"
	"github.com/hollson/dbcoder/utils"
	_ "github.com/lib/pq"
)

// PostgreSQL数据库代码生成器
type Gen struct {
	Source  string   // 连接字符串
	ignores []string // 忽略的表
}

func New(host string, port int, user, auth, dbname string, ignores []string) *Gen {
	if port == 0 {
		port = 5432
	}
	if len(user) == 0 {
		user = "postgres"
	}
	if len(auth) == 0 {
		auth = "postgres"
	}
	source := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, auth, host, port, dbname)
	return &Gen{source, ignores}
}

// 查询数据库表清单SQL
func (g *Gen) tablesSQL() string {
	return `SELECT a.tablename,
			COALESCE(c.description,'') AS comment
			FROM pg_tables a
			LEFT JOIN pg_class b on a.tablename=b.relname
			LEFT JOIN pg_description c on  b.oid=c.objoid and c.objsubid=0
			WHERE a.schemaname='public';`
}

// 查询数据表定义SQL
func (g *Gen) columnsSQL(tableName string) string {
	return fmt.Sprintf(`
SELECT a.attname                                       AS field_name,       
       a.attlen                                        AS field_size,
       a.atttypid::regtype                             AS field_type,
       COALESCE(ct.contype = 'p', FALSE)               AS is_primary_key,   
       a.attnotnull                                    AS not_null,         
       COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,    
       COALESCE(b.description, '')                     AS comment           
FROM pg_attribute a
         INNER JOIN ONLY pg_class C ON C.oid = a.attrelid
         INNER JOIN ONLY pg_namespace n ON n.oid = C.relnamespace
         LEFT JOIN pg_constraint ct ON ct.conrelid = C.oid AND a.attnum = ANY (ct.conkey) AND ct.contype = 'p'
         LEFT JOIN pg_attrdef ad ON ad.adrelid = C.oid AND ad.adnum = a.attnum
         LEFT JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid
         left join pg_type t on a.atttypid = t.oid
WHERE a.attisdropped = FALSE
  AND a.attnum > 0
  AND n.nspname = 'public'
  AND C.relname = '%s'
ORDER BY a.attnum;
`, tableName)
}

func (g *Gen) Tables() (ret []schema.Table, err error) {
	_db, err := sql.Open("postgres", g.Source)
	if err != nil {
		return nil, err
	}
	fmt.Printf(" 🛢 连接数据库: %s\n", g.Source)
	defer _db.Close()

	rows, err := _db.Query(g.tablesSQL())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t = schema.Table{}
		if err := rows.Scan(&t.Name, &t.Comment); err != nil {
			return nil, err
		}
		if utils.MatchAny(t.Name, g.ignores...) {
			continue
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

func (g *Gen) columns(tableName string, db *sql.DB) (ret []schema.Column, err error) {
	rows, err := db.Query(g.columnsSQL(tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var t = schema.Column{}
	for rows.Next() {
		if err := rows.Scan(
			&t.Name,
			&t.Size,
			&t.Type,
			&t.Primary,
			&t.NotNull,
			&t.Default,
			&t.Comment); err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}
	return
}

// Postgresql类型映射的Golang数据类型
//  参考：http://www.postgres.cn/docs/12/
//       http://www.postgres.cn/docs/12/datatype.html
func (g *Gen) TypeMapping(_type string) schema.Type {
	// 数组
	if utils.HasAny(_type, "[]") {
		switch {
		// 布尔：
		case utils.HasAny(_type, "boolean"):
			return schema.Type{Name: "pq.BoolArray", Pack: "github.com/lib/pq"}
			// 字节数组
		case utils.HasAny(_type, "bytea"):
			return schema.Type{Name: "pq.ByteaArray", Pack: "github.com/lib/pq"}
			// 整数：
		case utils.HasAny(_type, "smallint", "integer"):
			return schema.Type{Name: "pq.Int32Array", Pack: "github.com/lib/pq"}
		case utils.HasAny(_type, "bigint", "timestamp"):
			return schema.Type{Name: "pq.Int64Array", Pack: "github.com/lib/pq"}
			// 浮点数：
		case utils.HasAny(_type, "real"):
			return schema.Type{Name: "pq.Float32Array", Pack: "github.com/lib/pq"} // 单精度
		case utils.HasAny(_type, "double"):
			return schema.Type{Name: "pq.Float64Array", Pack: "github.com/lib/pq"} // 双精度
		case utils.HasAny(_type, "numeric", "decimal", "money"):
			return schema.Type{Name: "[]decimal.Decimal", Pack: "github.com/shopspring/decimal"}
			// 字符串：
		case utils.HasAny(_type, "uuid", "text", "character", "cidr", "inet", "macaddr", "interval"):
			return schema.Type{Name: "pq.StringArray", Pack: "github.com/lib/pq"}
			// 时间：
		case utils.HasAny(_type, "time with", "date"):
			return schema.Type{Name: "[]time.Time", Pack: "time"}
		default:
			return schema.Type{Name: "pq.GenericArray", Pack: "github.com/lib/pq"}
		}
	}

	// 标量数据
	switch {
	// 布尔：
	case utils.HasAny(_type, "boolean"):
		return schema.Type{Name: "bool"}
		// 字节数组
	case utils.HasAny(_type, "bytea"):
		return schema.Type{Name: "[]byte"}
		// 整数：
	case utils.HasAny(_type, "smallint"):
		return schema.Type{Name: "int16"}
	case utils.HasAny(_type, "integer"):
		return schema.Type{Name: "int32"}
	case utils.HasAny(_type, "bigint", "timestamp"):
		return schema.Type{Name: "int64"}
		// 浮点数：
	case utils.HasAny(_type, "real"):
		return schema.Type{Name: "float32"} // 单精度
	case utils.HasAny(_type, "double"):
		return schema.Type{Name: "float64"} // 双精度
	case utils.HasAny(_type, "numeric", "decimal", "money"):
		return schema.Type{Name: "decimal.Decimal", Pack: "github.com/shopspring/decimal"}
		// 字符串：
	case utils.HasAny(_type, "uuid", "text", "character", "cidr", "inet", "macaddr", "interval"):
		return schema.Type{Name: "string"}
		// 时间：
	case utils.HasAny(_type, "time with", "date"):
		return schema.Type{Name: "time.Time", Pack: "time"}
	default:
		return schema.Type{Name: "interface{}"}
	}
}
