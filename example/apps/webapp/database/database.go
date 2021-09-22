package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/razorpay/devstack/example/apps/webapp/pkg/helpers"
	"github.com/razorpay/devstack/example/apps/webapp/pkg/tracing"
)

//DB ...
type DB struct {
	gorm.DB
}

//DBConnector ...
var dbConnector *gorm.DB
var DBTracer = &tracing.JaegerTracer{
	Tracer: nil,
	Closer: nil,
}

func Init() {
	dbConnector, _ = connectDB()
}

func connectDB() (*gorm.DB, error) {
	basePath, err := os.Getwd()
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to get Base Path for GormDB"))
		return nil, err
	}
	dbPath := fmt.Sprintf("%s/gorm.db", basePath)
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Panic(fmt.Sprintf("Database Error:%s", err))
	}

	log.Print("Connected to DB")

	dbServiceName := tracing.DBPeerService
	log.Print("Initializing DB tracer: %s", dbServiceName)
	DBTracer = tracing.InitTracing(dbServiceName)

	helpers.AddGormCallbacks(db, DBTracer)
	//AddGormCallbacks(db)
	//db.Callback().Create().Replace("gorm:update_time_stamp", setCreatedTimeStamp)
	//db.Callback().Update().Replace("gorm:update_time_stamp", setUpdatedTimeStamp)

	// Enable connecting pooling
	//db.DB().SetMaxIdleConns(MaxIdleConnections)
	//db.DB().SetMaxOpenConns(MaxOpenConnections)
	return db, err
}

//GetDB ...
func GetDB(ctx context.Context) *gorm.DB {
	if ctx != nil {
		return helpers.SetSpanToGorm(ctx, dbConnector)
	}
	return dbConnector
}

// GetDBTracer
func GetDBTracer() *tracing.JaegerTracer {
	return DBTracer
}
