package query

import (
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const gormSpanKey = "__gorm_span"
const gormTime = "__gorm_time"

const (
	callBackBeforeName = "opentracing:before"
	callBackAfterName  = "opentracing:after"
)

var _ gorm.Plugin = (*OpentracingPlugin)(nil)

type OpentracingPlugin struct{}

// NewOpentracingPlugin 创建一个opentracing插件
func NewOpentracingPlugin() *OpentracingPlugin {
	return &OpentracingPlugin{}
}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	if err = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before); err != nil {
		return err
	}
	if err = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before); err != nil {
		return err
	}

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	if err = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after); err != nil {
		return err
	}
	if err = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after); err != nil {
		return err
	}
	return
}

func before(db *gorm.DB) {
	_, span := otel.Tracer("gorm").Start(db.Statement.Context, "gorm")
	// 利用db实例去传递span
	db.InstanceSet(gormSpanKey, span)
	db.InstanceSet(gormTime, time.Now())
}

func after(db *gorm.DB) {
	_span, isExist := db.InstanceGet(gormSpanKey)
	if !isExist {
		return
	}

	span, ok := _span.(oteltrace.Span)
	if !ok {
		return
	}
	defer span.End()

	// Error
	if db.Error != nil {
		span.SetAttributes(attribute.String("error", db.Error.Error()))
	}
	// sql
	span.SetAttributes(attribute.String("sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))
	// rows
	span.SetAttributes(attribute.Int64("rows", db.RowsAffected))
	// elapsed
	_time, isExist := db.InstanceGet(gormTime)
	if !isExist {
		return
	}
	startTime, ok := _time.(time.Time)
	if ok {
		span.SetAttributes(attribute.String("elapsed", time.Since(startTime).String()))
	}
}
