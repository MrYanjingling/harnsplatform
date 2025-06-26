package data

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"harnsplatform/internal/biz"
	"harnsplatform/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/mattn/go-sqlite3"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewThingTypesRepo)

// Data .
type Data struct {
	DB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, log *log.Helper) (*Data, func(), error) {
	data := &Data{}

	driver := c.GetDatabase().GetDriver()
	switch driver {
	case "sqlite":
		if db, err := gorm.Open(sqlite.Open(c.GetDatabase().GetSource()), &gorm.Config{}); err != nil {
			return nil, nil, err
		} else {
			data.DB = db
		}
	case "mysql":
		if db, err := gorm.Open(mysql.Open(c.GetDatabase().GetSource()), &gorm.Config{}); err != nil {
			return nil, nil, err
		} else {
			data.DB = db
		}
	}

	sqlDB, err := data.DB.DB()
	if err != nil {
		return nil, nil, err
	}
	// 配置连接池
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // 连接最大存活时间

	if err := data.DB.AutoMigrate(&biz.ThingTypes{}, &biz.Things{}); err != nil {
		log.Fatalf("failed to migrate thingTypes model: %v", err)
	}

	cleanup := func() {
		log.Info("closing the data resources")
		if db, err := data.DB.DB(); err != nil || db.Close() != nil {
			log.Error("failed to close datasource")
		}
	}
	return data, cleanup, nil
}
