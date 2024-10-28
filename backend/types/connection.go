package types

import "github.com/google/uuid"

type Connection struct {
	ID     string          `json:"id" yaml:"id"`
	Name   string          `json:"name" yaml:"name"`
	Config MysqlConnection `json:"config" yaml:"config"`
}

type Connections []Connection

// Init ID
func (c *Connection) InitID() {
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
}
