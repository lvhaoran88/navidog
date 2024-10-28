package services

import (
	"context"
	"sync"
)

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
