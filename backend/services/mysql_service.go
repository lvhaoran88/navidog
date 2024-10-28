package services

import (
	"context"
	"database/sql"
	"navidog/backend/types"
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
