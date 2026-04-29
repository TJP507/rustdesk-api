package orm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type MysqlConfig struct {
	Dsn          string
	MaxIdleConns int
	MaxOpenConns int
}

func NewMysql(mysqlConf *MysqlConfig, logwriter logger.Writer) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               mysqlConf.Dsn, // DSN data source name
		DefaultStringSize: 256,           // Default length for string type fields
		//DisableDatetimePrecision:  true,                    // Disable datetime precision; not supported by MySQL databases prior to 5.6
		//DontSupportRenameIndex:    true,                    // Rename indexes by dropping and re-creating them; not supported by MySQL prior to 5.7 or MariaDB
		//DontSupportRenameColumn:   true,                    // Use `change` to rename columns; not supported by MySQL prior to 8 or MariaDB
		//SkipInitializeWithVersion: false,                   // Auto-configure based on the current MySQL version
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(
			logwriter, // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Warn, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL log
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		fmt.Println(err)
	}
	sqlDB, err2 := db.DB()
	if err2 != nil {
		fmt.Println(err2)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open database connections
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConns)

	return db
}
