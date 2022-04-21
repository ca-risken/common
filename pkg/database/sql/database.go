package sql

import (
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Open(dsn string, logMode bool, maxConn int) (*gorm.DB, error) {
	sqltrace.Register("mysql", &mysqldriver.MySQLDriver{})
	dbConn, err := sqltrace.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	dbConn.SetMaxOpenConns(maxConn)
	dbConn.SetConnMaxLifetime(time.Duration(maxConn/2) * time.Second)

	db, err := gormtrace.Open(mysql.New(mysql.Config{
		Conn: dbConn,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return nil, err
	}
	if logMode {
		db.Logger.LogMode(logger.Info)
	}
	return db, nil
}
