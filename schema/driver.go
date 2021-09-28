// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package schema

// database driver type
type Driver int

const (
	Unknown Driver = iota
	MySQL
	PostgreSQL
	SQLite
	MongoDB
	MariaDB
	Oracle
	SQLServer
)

// driver name
func (d Driver) String() string {
	return []string{"Unknown", "MySQL", "PostgreSQL", "SQLite", "MongoDB", "MariaDB", "Oracle", "SQLServer"}[d]
}

// driver short name
func (d Driver) Name() string {
	return []string{"unknown", "mysql", "pgsql", "sqlite", "mongo", "maria", "oracle", "mssql"}[d]
}

// get driver value by driver short name
func DriverValue(name string) Driver {
	if d, ok := map[string]Driver{
		MySQL.Name():      MySQL,
		PostgreSQL.Name(): PostgreSQL,
		SQLite.Name():     SQLite,
		MongoDB.Name():    MongoDB,
		MariaDB.Name():    MariaDB,
		Oracle.Name():     Oracle,
		MySQL.Name():      SQLServer,
	}[name]; ok {
		return d
	}
	return Unknown
}
