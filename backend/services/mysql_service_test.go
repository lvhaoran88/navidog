package services

import (
	"navidog/backend/types"
	"testing"
)

func TestTransferData(t *testing.T) {
	source := types.MysqlConnection{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "root",
		Database: "cmdb",
	}
	target := types.MysqlConnection{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "root",
		Database: "test_cmdb",
	}

	mysqlService := MysqlService()
	mysqlService.TransferData(source, target, []string{"auth_group", "auth_group_permissions"})
}
