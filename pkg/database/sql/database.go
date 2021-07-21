package sql

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Open(dsn string, logMode bool) (*gorm.DB, error) {
	instrumentedDB, err := xray.SQLContext("mysql", dsn)
	if err != nil {
		return nil, err
	}
	/*
	 In initializing the context.Context which doesn't have segment is used to create Subsegment and some error occur.
	 To avoid shutdown in it, Must set AWS_XRAY_CONTEXT_MISSING=LOG_ERROR.

	 * the point of error
	 https://github.com/go-gorm/gorm/blob/v1.21.11/gorm.go#L195
	 https://github.com/go-gorm/mysql/blob/v1.1.1/mysql.go#L89

	 * error logs
	 ------
	 [ERROR] Suppressing AWS X-Ray context missing panic: failed to begin subsegment named 'mimosa': segment cannot be found.
	 [ERROR] Suppressing AWS X-Ray context missing panic: failed to get segment from context since segment is nil
	 [ERROR] Suppressing AWS X-Ray context missing panic: failed to end subsegment: subsegment 'mimosa' cannot be found.
	 ------
	*/
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: instrumentedDB,
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return nil, err
	}
	if logMode {
		db.Logger.LogMode(logger.Info)
	}
	return db, nil
}
