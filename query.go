package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type (
	Tracer interface {
		IsEnableTrace() bool
		// OpenTrace 开启trace
		OpenTrace() Tracer
		// CloseTrace 关闭trace
		CloseTrace() Tracer
	}

	ICtx interface {
		GetCtx() context.Context
	}

	// IAction 操作接口
	IAction[T any] interface {
		IOperation[T]
		IOperationX[T]
		IBind[T]

		IAssociation
	}

	action[T any] struct {
		db    *gorm.DB
		ctx   context.Context
		table schema.Tabler

		IAssociation
		IOperation[T]
		IOperationX[T]
		Tracer
	}

	ActionOption[T any] func(a *action[T])
)

// NewAction 创建GORM操作接口实例
func NewAction[T any](opts ...ActionOption[T]) IAction[T] {
	ac := action[T]{
		ctx:         context.Background(),
		IOperation:  NewOperationImpl[T](),
		IOperationX: NewOperationXImpl[T](),
		Tracer:      NewITracer(),
	}

	for _, opt := range opts {
		opt(&ac)
	}

	if ac.table != nil {
		ac.db = ac.db.Table(ac.table.TableName())
	}

	if ac.IAssociation == nil {
		ac.IAssociation = NewDefaultAssociation(ac.db)
	}

	return &ac
}

// DB 获取DB, 包含了Table或Model, 用于链式操作
func (a *action[T]) DB() *gorm.DB {
	if a.table != nil {
		return a.db.Table(a.table.TableName())
	}
	var m T
	return a.db.Model(&m)
}

// Clauses 设置Clauses
func (a *action[T]) Clauses(condList ...clause.Expression) IAction[T] {
	a.db = a.db.Clauses(condList...)
	return a
}

// Order 跳转到排序动作
func (a *action[T]) Order(column string) IOrder[T] {
	return NewOrder[T](column).WithIAction(a)
}

// WithDB 设置DB, 一般用于事务, 这里使用事务的DB, 也可以设置新的DB用于链式操作
func (a *action[T]) WithDB(db *gorm.DB) IAction[T] {
	a.db = db
	return a
}

// WithContext 设置上下文Ctx
func (a *action[T]) WithContext(ctx context.Context) IAction[T] {
	a.ctx = ctx
	return a
}

// WithTable 设置Table, 这里传递的是实现了schema.Tabler接口的结构体
func (a *action[T]) WithTable(tabler schema.Tabler) IAction[T] {
	a.table = tabler
	return a
}

// WithModel 设置Model, Model规范参考: https://gorm.io/zh_CN/docs/models.html
func (a *action[T]) WithModel(model any) IAction[T] {
	a.db = a.db.Model(model)
	return a
}

// Preload 预加载, 参考: https://gorm.io/zh_CN/docs/preload.html
func (a *action[T]) Preload(preloadKey string, wheres ...ScopeMethod) IAction[T] {
	a.db = a.db.Preload(preloadKey, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(wheres...)
	})
	return a
}

// Joins 设置关联, 参考: https://gorm.io/zh_CN/docs/preload.html#Joins-%E9%A2%84%E5%8A%A0%E8%BD%BD
func (a *action[T]) Joins(joinsKey string, wheres ...ScopeMethod) IAction[T] {
	a.db = a.db.Joins(joinsKey, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(wheres...)
	})
	return a
}

// Scopes 设置作用域, 参考: https://gorm.io/zh_CN/docs/scopes.html
func (a *action[T]) Scopes(wheres ...ScopeMethod) IAction[T] {
	a.db = a.db.Scopes(wheres...)
	return a
}
