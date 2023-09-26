package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var _ IAction[any] = (*Action[any])(nil)

type IOpration[T any] interface {
	// First 查询单条数据
	First(wheres ...Scopemethod) (*T, error)
	// Last 查询单条数据
	Last(wheres ...Scopemethod) (*T, error)
	// List 查询多条数据
	List(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error)
	// Count 查询数量
	Count(wheres ...Scopemethod) (int64, error)
	// Create 创建数据
	Create(m *T) error
	// Update 更新数据
	Update(m *T, wheres ...Scopemethod) error
	// Delete 删除数据
	Delete(wheres ...Scopemethod) error
	// ForcedDelete 强制删除数据
	ForcedDelete(wheres ...Scopemethod) error
}

type IBind[T any] interface {
	// WithDB 设置DB
	WithDB(db *gorm.DB) IAction[T]
	// WithContext 设置Ctx
	WithContext(ctx context.Context) IAction[T]
	// WithTabler 设置Tabler
	WithTabler(tabler schema.Tabler) IAction[T]
	// WithModel 设置Model
	WithModel(model any) IAction[T]
	// Preload 预加载
	Preload(preloadKey string, wheres ...Scopemethod) IAction[T]
	// Joins 设置关联
	Joins(joinsKey string, wheres ...Scopemethod) IAction[T]
	// Scopes 设置作用域
	Scopes(wheres ...Scopemethod) IAction[T]

	// Order 排序
	Order(column string) IOrder[T]

	// Clauses 设置Clauses
	Clauses(conds ...clause.Expression) IAction[T]
}

type IAction[T any] interface {
	IOpration[T]
	IBind[T]

	DB() *gorm.DB
}

type Action[T any] struct {
	db  *gorm.DB
	ctx context.Context
	schema.Tabler
}

type ActionOption[T any] func(a *Action[T])

func NewAction[T any](opts ...ActionOption[T]) *Action[T] {
	action := Action[T]{
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(&action)
	}

	if action.Tabler != nil {
		action.db = action.db.Table(action.Tabler.TableName())
	}

	return &action
}

func WithDB[T any](db *gorm.DB) ActionOption[T] {
	return func(a *Action[T]) {
		a.db = db
	}
}

func WithContext[T any](ctx context.Context) ActionOption[T] {
	return func(a *Action[T]) {
		a.ctx = ctx
	}
}

func WithTabler[T any](tabler schema.Tabler) ActionOption[T] {
	return func(a *Action[T]) {
		a.Tabler = tabler
	}
}

func (a *Action[T]) DB() *gorm.DB {
	if a.Tabler != nil {
		return a.db.Table(a.Tabler.TableName())
	}
	var m T
	return a.db.Model(&m)
}

func (a *Action[T]) Clauses(conds ...clause.Expression) IAction[T] {
	a.db = a.db.Clauses(conds...)
	return a
}

func (a *Action[T]) Order(column string) IOrder[T] {
	return NewOrder[T](column).WithIAction(a)
}

func (a *Action[T]) WithDB(db *gorm.DB) IAction[T] {
	a.db = db
	return a
}

func (a *Action[T]) WithContext(ctx context.Context) IAction[T] {
	a.ctx = ctx
	return a
}

func (a *Action[T]) WithTabler(tabler schema.Tabler) IAction[T] {
	a.Tabler = tabler
	return a
}

func (a *Action[T]) WithModel(model any) IAction[T] {
	a.db = a.db.Model(model)
	return a
}

func (a *Action[T]) Preload(preloadKey string, wheres ...Scopemethod) IAction[T] {
	a.db = a.db.Preload(preloadKey, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(wheres...)
	})
	return a
}

func (a *Action[T]) Joins(joinsKey string, wheres ...Scopemethod) IAction[T] {
	a.db = a.db.Joins(joinsKey, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(wheres...)
	})
	return a
}

func (a *Action[T]) Scopes(wheres ...Scopemethod) IAction[T] {
	a.db = a.db.Scopes(wheres...)
	return a
}

// First 查询单条数据
func (a *Action[T]) First(wheres ...Scopemethod) (*T, error) {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).First(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// Last 查询单条数据
func (a *Action[T]) Last(wheres ...Scopemethod) (*T, error) {
	var m T

	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Last(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// List 查询多条数据
func (a *Action[T]) List(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error) {
	var ms []*T

	db := a.DB().WithContext(a.ctx).Scopes(wheres...)
	if pgInfo != nil {
		var total int64
		if err := db.WithContext(a.ctx).Count(&total).Error; err != nil {
			return nil, err
		}
		pgInfo.SetTotal(total)
		db = db.Scopes(Paginate(pgInfo))
	}

	if err := db.Find(&ms).Error; err != nil {
		return nil, err
	}

	return ms, nil
}

// Count 查询数量
func (a *Action[T]) Count(wheres ...Scopemethod) (int64, error) {
	var total int64

	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Create 创建数据
func (a *Action[T]) Create(newModel *T) error {
	return a.DB().WithContext(a.ctx).Create(newModel).Error
}

// Update 更新数据
func (a *Action[T]) Update(newModel *T, wheres ...Scopemethod) error {
	return a.DB().WithContext(a.ctx).Scopes(wheres...).Updates(newModel).Error
}

// Delete 删除数据
func (a *Action[T]) Delete(wheres ...Scopemethod) error {
	var m T
	return a.DB().WithContext(a.ctx).Scopes(wheres...).Delete(&m).Error
}

// ForcedDelete 强制删除数据
func (a *Action[T]) ForcedDelete(wheres ...Scopemethod) error {
	var m T
	return a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).Delete(&m).Error
}
