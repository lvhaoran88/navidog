package services

import (
	"context"
	"database/sql"
	"fmt"
	"navidog/backend/types"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

const mysqlDriverName = "mysql"

type mysqlService struct {
	ctx context.Context
}

var mysqlServiceInstance *mysqlService
var mysqlServiceInstanceOnce sync.Once

func MysqlService() *mysqlService {
	if mysqlServiceInstance == nil {

		mysqlServiceInstanceOnce.Do(func() {
			mysqlServiceInstance = &mysqlService{
				ctx: context.Background(),
			}
		})
	}
	return mysqlServiceInstance
}

func (ms *mysqlService) Startup(ctx context.Context) {
	ms.ctx = ctx
}

// TestConnection
func (ms *mysqlService) TestConnection(connection types.MysqlConnection) (resp types.JSResp) {
	db, err := sql.Open(mysqlDriverName, connection.FormatDSN())
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	return
}

// ListDatabases
func (ms *mysqlService) ListDatabases(connection types.MysqlConnection) (resp types.JSResp) {
	db, err := sql.Open(mysqlDriverName, connection.FormatDSN())
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer db.Close()

	var databases []string
	if err = db.QueryRow("SHOW DATABASES").Scan(&databases); err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	resp.Data = databases
	return
}

// ListTables
func (ms *mysqlService) ListTables(connection types.MysqlConnection) (resp types.JSResp) {
	db, err := sql.Open(mysqlDriverName, connection.FormatDSN())
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer db.Close()

	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			resp.Message = err.Error()
			return
		}
		tables = append(tables, table)
	}
	resp.Success = true
	resp.Data = tables
	return
}

// TransferData 传输数据 将源数据库指定表的数据传输到目标数据库对应的表
func (ms *mysqlService) TransferData(source types.MysqlConnection, target types.MysqlConnection, tables []string) (resp types.JSResp) {

	sourceDb, err := sql.Open(mysqlDriverName, source.FormatDSN())
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer sourceDb.Close()

	targetDb, err := sql.Open(mysqlDriverName, target.FormatDSN())
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer targetDb.Close()

	for _, table := range tables {
		// 删除目标数据库中对应的表
		// 取消外键检测
		if _, err := targetDb.Exec("SET FOREIGN_KEY_CHECKS = 0;"); err != nil {
			resp.Message = err.Error()
			return
		}
		if _, err = targetDb.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
			resp.Message = err.Error()
			return
		}
		// 恢复外键检测
		if _, err = targetDb.Exec("SET FOREIGN_KEY_CHECKS = 1;"); err != nil {
			resp.Message = err.Error()
			return
		}

		// 查询源数据库中对应表结构
		var createTableSql string
		if err := sourceDb.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", table)).Scan(&table, &createTableSql); err != nil {
			resp.Message = err.Error()
			return
		}
		// 创建目标数据库中对应的表
		if _, err := targetDb.Exec(createTableSql); err != nil {
			resp.Message = err.Error()
			return
		}

		// 查询源数据库中对应表数据
		rows, err := sourceDb.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			resp.Message = err.Error()
			return
		}
		defer rows.Close()

		// 获取列名
		columns, err := rows.Columns()
		if err != nil {
			resp.Message = err.Error()
			return
		}

		// 读取数据
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		for rows.Next() {
			// 后面使用 channel + goroutine 来提高效率
			if err := rows.Scan(valuePtrs...); err != nil {
				resp.Message = err.Error()
				return
			}

			// 插入数据到目标数据库中对应的表
			if _, err := targetDb.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Repeat("?,", len(columns)-1)+"?"), values...); err != nil {
				resp.Message = err.Error()
				return
			}
		}

	}
	resp.Success = true
	return

}
