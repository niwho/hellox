package datasource

import (
	"github.com/niwho/hellox/common"
	"github.com/niwho/hellox/logs"
)

var DataSourceIns *DataSource

type DataSource struct {
	DbSource *common.DBClient
}

func InitDataSource() {
	DataSourceIns = &DataSource{
		DbSource: NewSqliteDBClient("hellox.db"),
	}
}
