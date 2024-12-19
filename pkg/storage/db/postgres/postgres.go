package postgres

import (
	"errors"
	"fmt"
	"ginboot/pkg/conf"
	"ginboot/pkg/logger"
	"ginboot/pkg/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *conf.PostgresConfig) error {
	if cfg == nil {
		return errors.New("postgres config is nil")
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.Ip, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.Ssl)

	var err error
	storage.DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: NewPostgresLogger(),
	})

	if err != nil {
		return err
	}

	logger.Out.Debug("postgres init finished")
	return nil
}
