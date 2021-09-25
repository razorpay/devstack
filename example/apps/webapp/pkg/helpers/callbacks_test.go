package helpers

import (
	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"testing"
)

/**
Taken from https://github.com/smacker/opentracing-gorm/blob/master/otgorm.go.
This test case is more of including the test coverage,
also keeping an eye for any changes in the callback.
*/

var tracer *mocktracer.MockTracer
var gDB *gorm.DB

func init() {
	gDB = initDB()
	tracer = mocktracer.New()
	opentracing.SetGlobalTracer(tracer)
}

type Product struct {
	gorm.Model
	Code string
}

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Product{})
	db.Create(&Product{Code: "L1212"})
	AddGormCallbacks(db)
	return db
}

func Handler(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handler")
	defer span.Finish()

	db := SetSpanToGorm(ctx, gDB)

	var product Product
	db.First(&product, 1)
}

func TestPool(t *testing.T) {
	Handler(context.Background())
	spans := tracer.FinishedSpans()
	if len(spans) != 2 {
		t.Fatalf("should be 2 finished spans but there are %d: %v", len(spans), spans)
	}

	sqlSpan := spans[0]
	if sqlSpan.OperationName != "sql" {
		t.Errorf("first span operation should be sql but it's '%s'", sqlSpan.OperationName)
	}

	expectedTags := map[string]interface{}{
		"error":        false,
		"db.table":     "products",
		"db.method":    "SELECT",
		"db.type":      "sql",
		"db.statement": `SELECT * FROM "products"  WHERE "products"."deleted_at" IS NULL AND (("products"."id" = 1)) ORDER BY "products"."id" ASC LIMIT 1`,
		"db.err":       false,
		"db.count":     int64(1),
		"peer.service": "",
		"peer.address": "",
		"peer.port":    uint16(0),
		"span.kind":    ext.SpanKindEnum("client"),
	}

	sqlTags := sqlSpan.Tags()

	if len(sqlTags) != len(expectedTags) {
		t.Errorf("sql span should have %d tags but it has %d", len(expectedTags), len(sqlTags))
	}

	for name, expected := range expectedTags {
		value, ok := sqlTags[name]
		if !ok {
			t.Errorf("sql span doesn't have tag '%s'", name)
			continue
		}
		if value != expected {
			t.Errorf("sql span tag '%s' should have value '%s' but it has '%s'", name, expected, value)
		}
	}

	if spans[1].OperationName != "handler" {
		t.Errorf("second span operation should be handler but it's '%s'", spans[1].OperationName)
	}
}
