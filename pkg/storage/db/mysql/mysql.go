package mysql

import (
	"server/pkg/conf"
)

func Init(cfg *conf.MysqlConfig) (err error) {
	if cfg == nil {
		return
	}
	return
}
