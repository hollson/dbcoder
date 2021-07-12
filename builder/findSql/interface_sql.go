package findSql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/hollson/dbcoder/dialect/mysql"
	"github.com/hollson/dbcoder/dialect/pgsql"
)

type SqlInterface interface {
	Init(userName, pwd, host string, port int, dbName string) error
	FindTableString() string
	FindColumnsString() string
	DBNameString() string
	GetDB() *sql.DB
}

var driverMap map[string]SqlInterface

func init() {
	driverMap = make(map[string]SqlInterface, 0)
	driverMap["mysql"] = &mysql.MySql{}
	driverMap["pgsql"] = &pgsql.PostgresSql{}
}

func SelectDriver(driver string) (SqlInterface, error) {
	driver = strings.ToLower(driver)
	model, ok := driverMap[driver]
	if !ok {
		fmt.Println(driver)
		return nil, errors.New("错误：该程序不包含该数据库的导出功能，请检查 -driver 是否设置正确...")
	}
	return model, nil
}
