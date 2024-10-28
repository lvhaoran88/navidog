package services

import (
	"context"
	"sync"
)

type connectionService struct {
	ctx context.Context
}

var connectionServiceInstance *connectionService
var connectionServiceInstanceOnce sync.Once

// ConnectionService
func ConnectionService() *connectionService {
	if connectionServiceInstance == nil {
		connectionServiceInstanceOnce.Do(func() {
			connectionServiceInstance = &connectionService{}
		})
	}
	return connectionServiceInstance
}

// Startup
func (cs *connectionService) Startup(ctx context.Context) {
	cs.ctx = ctx
}
