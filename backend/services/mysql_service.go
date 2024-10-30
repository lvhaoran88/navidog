package services

import (
	"context"
	"database/sql"
	"fmt"
	"navidog/backend/types"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

// Version 获取数据库版本
func (ms *mysqlService) Version(connection types.MysqlConnection) (resp types.JSResp) {

	db, err := sql.Open(mysqlDriverName, connection.FormatDSN())
	if err != nil {
		resp.Message = err.Error()
		return
	}
	defer db.Close()

	var version string
	if err := db.QueryRow("SELECT VERSION()").Scan(&version); err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Success = true
	resp.Data = version
	return
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
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		resp.Message = err.Error()
		return
	}
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			resp.Message = err.Error()
			return
		}
		databases = append(databases, database)
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
	defer func() {
		if err := recover(); err != nil {
			resp.Message = fmt.Sprintf("传输数据失败: %v", err)
			runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("传输数据失败: %v", err))
		} else {
			resp.Success = true
		}
	}()

	runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("连接数据库: %s:%d", source.Host, source.Port))
	sourceDb, err := sql.Open(mysqlDriverName, source.FormatDSN())
	if err != nil {
		panic(err)
	}
	defer sourceDb.Close()
	runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("连接数据库: %s:%d 成功", source.Host, source.Port))

	runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("连接数据库: %s:%d", target.Host, target.Port))
	targetDb, err := sql.Open(mysqlDriverName, target.FormatDSN())
	if err != nil {
		panic(err)
	}
	defer targetDb.Close()
	runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("连接数据库: %s:%d 成功", target.Host, target.Port))
	runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("准备传输数据, 共 %d 张表", len(tables)))
	for index, table := range tables {
		runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("开始传输第 %d 张表: %s", index+1, table))
		// 删除目标数据库中对应的表
		// 取消外键检测
		if _, err := targetDb.Exec("SET FOREIGN_KEY_CHECKS = 0;"); err != nil {
			panic(err)
		}
		if _, err = targetDb.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
			panic(err)
		}

		// 查询源数据库中对应表结构
		var createTableSql string
		if err := sourceDb.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", table)).Scan(&table, &createTableSql); err != nil {
			panic(err)
		}
		// 创建目标数据库中对应的表
		if _, err := targetDb.Exec(createTableSql); err != nil {
			panic(err)
		}

		// 查询源数据库中对应表数据
		rows, err := sourceDb.Query(fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		// 获取列名
		columns, err := rows.Columns()
		if err != nil {
			panic(err)
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
				panic(err)
			}

			// 插入数据到目标数据库中对应的表
			if _, err := targetDb.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Repeat("?,", len(columns)-1)+"?"), values...); err != nil {
				panic(err)
			}

		}
		// 恢复外键检测
		if _, err = targetDb.Exec("SET FOREIGN_KEY_CHECKS = 1;"); err != nil {
			panic(err)
		}

		runtime.EventsEmit(ms.ctx, "transferDataProgress", index+1, len(tables))
		runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("传输第 %d 张表: %s 完成", index+1, table))

	}
	runtime.EventsEmit(ms.ctx, "transferDataLog", fmt.Sprintf("传输数据完成, 共 %d 张表", len(tables)))

	return
}
