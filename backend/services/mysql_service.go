package services

import (
	"context"
	"database/sql"
	"navidog/backend/types"
	"sync"
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
