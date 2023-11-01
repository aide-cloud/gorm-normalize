package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// IBind 绑定操作, 用于链式操作
type IBind[T any] interface {
	// WithDB 设置DB
	WithDB(db *gorm.DB) IAction[T]
	// WithContext 设置Ctx
	WithContext(ctx context.Context) IAction[T]
	// WithTable 设置Table
	WithTable(table schema.Tabler) IAction[T]
	// WithModel 设置Model
	WithModel(model any) IAction[T]
	// Preload 预加载
	Preload(preloadKey string, wheres ...ScopeMethod) IAction[T]
	// Joins 设置关联
	Joins(joinsKey string, wheres ...ScopeMethod) IAction[T]
	// Scopes 设置作用域
	Scopes(wheres ...ScopeMethod) IAction[T]

	// Order 排序
	Order(column string) IOrder[T]

	// Clauses 设置Clauses
	Clauses(condList ...clause.Expression) IAction[T]

	// OpenTrace 开启trace
	OpenTrace() IAction[T]
	// CloseTrace 关闭trace
	CloseTrace() IAction[T]
}
