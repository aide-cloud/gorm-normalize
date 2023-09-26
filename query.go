package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _ IAction[any] = (*Action[any])(nil)

type IAction[T any] interface {
	// First 查询单条数据
	First(wheres ...Scopemethod) (*T, error)
	// Last 查询单条数据
	Last(wheres ...Scopemethod) (*T, error)
	// List 查询多条数据
	List(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error)
	// Create 创建数据
	Create(m *T) error
	// Update 更新数据
	Update(m *T, wheres ...Scopemethod) error
	// Delete 删除数据
	Delete(wheres ...Scopemethod) error

	// WithDB 设置DB
	WithDB(db *gorm.DB) IAction[T]
	// WithContext 设置Ctx
	WithContext(ctx context.Context) IAction[T]
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

func (a *Action[T]) WithDB(db *gorm.DB) IAction[T] {
	a.db = db
	return a
}

func (a *Action[T]) WithContext(ctx context.Context) IAction[T] {
	a.ctx = ctx
	return a
}

// First 查询单条数据
func (a *Action[T]) First(wheres ...Scopemethod) (*T, error) {
	var m T
	if err := a.db.WithContext(a.ctx).Scopes(wheres...).First(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// Last 查询单条数据
func (a *Action[T]) Last(wheres ...Scopemethod) (*T, error) {
	var m T

	if err := a.db.WithContext(a.ctx).Scopes(wheres...).Last(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// List 查询多条数据
func (a *Action[T]) List(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error) {
	var ms []*T

	db := a.db.WithContext(a.ctx).Scopes(wheres...)

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

// Create 创建数据
func (a *Action[T]) Create(newModel *T) error {
	return a.db.WithContext(a.ctx).Create(newModel).Error
}

// Update 更新数据
func (a *Action[T]) Update(newModel *T, wheres ...Scopemethod) error {
	return a.db.WithContext(a.ctx).Scopes(wheres...).Updates(newModel).Error
}

// Delete 删除数据
func (a *Action[T]) Delete(wheres ...Scopemethod) error {
	var m T
	return a.db.WithContext(a.ctx).Scopes(wheres...).Delete(&m).Error
}
