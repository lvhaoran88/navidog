package services

import (
	"context"
	"navidog/backend/storage"
	"navidog/backend/types"
	"sync"

	"gopkg.in/yaml.v3"
)

type connectionService struct {
	ctx     context.Context
	storage *storage.ConnectionStorage
}

var connectionServiceInstance *connectionService
var connectionServiceInstanceOnce sync.Once

// ConnectionService
func ConnectionService() *connectionService {
	if connectionServiceInstance == nil {
		connectionServiceInstanceOnce.Do(func() {
			connectionServiceInstance = &connectionService{
				storage: storage.NewConnectionStorage(),
			}
		})
	}
	return connectionServiceInstance
}

// Startup
func (cs *connectionService) Startup(ctx context.Context) {
	cs.ctx = ctx
}

// defaultConnectios
func (cs *connectionService) defaultConnectios() types.Connections {
	return types.Connections{}
}

// listConnections
func (cs *connectionService) listConnections() types.Connections {
	b, err := cs.storage.Read()
	if err != nil {
		return cs.defaultConnectios()
	}
	var connections types.Connections
	if err = yaml.Unmarshal(b, &connections); err != nil {
		return cs.defaultConnectios()
	}
	return connections
}

// ListConnections
func (cs *connectionService) ListConnections() (resp types.JSResp) {
	resp.Data = cs.listConnections()
	resp.Success = true
	return
}

// SaveConnetion
func (cs *connectionService) SaveConnetion(connection types.Connection) (resp types.JSResp) {
	connections := cs.listConnections()
	if connection.ID == "" {
		// 新建场景
		connection.InitID()
		connections = append(connections, connection)
		b, err := yaml.Marshal(connections)
		if err != nil {
			resp.Message = err.Error()
			return
		}
		if err = cs.storage.Write(b); err != nil {
			resp.Message = err.Error()
			return
		}
	} else {
		//
		isMatch := false
		for i, c := range connections {
			if c.ID == connection.ID {
				connections[i] = connection
				isMatch = true
				break
			}
		}
		if !isMatch {
			connections = append(connections, connection)
		}
		b, err := yaml.Marshal(connections)
		if err != nil {
			resp.Message = err.Error()
			return
		}
		if err = cs.storage.Write(b); err != nil {
			resp.Message = err.Error()
			return
		}
	}
	resp.Success = true
	return
}
