package query

import (
	"context"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var _ IAction[any] = (*Action[any])(nil)

type IOpration[T any] interface {
	// First 查询单条数据
	First(wheres ...Scopemethod) (*T, error)
	// FirstWithTrashed 查询单条数据(包含软删除数据)
	FirstWithTrashed(wheres ...Scopemethod) (*T, error)
	// FirstByID 根据ID查询单条数据
	FirstByID(id uint) (*T, error)
	// FirstByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
	FirstByIDWithTrashed(id uint) (*T, error)
	// Last 查询单条数据
	Last(wheres ...Scopemethod) (*T, error)
	// LastWithTrashed 查询单条数据(包含软删除数据)
	LastWithTrashed(wheres ...Scopemethod) (*T, error)
	// LastByID 根据ID查询单条数据
	LastByID(id uint) (*T, error)
	// LastByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
	LastByIDWithTrashed(id uint) (*T, error)
	// List 查询多条数据
	List(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error)
	// ListWithTrashed 查询多条数据(包含软删除数据)
	ListWithTrashed(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error)
	// Count 查询数量
	Count(wheres ...Scopemethod) (int64, error)
	// CountWithTrashed 查询数量(包含软删除数据)
	CountWithTrashed(wheres ...Scopemethod) (int64, error)
	// Create 创建数据
	Create(m *T) error
	// Update 更新数据
	Update(m *T, wheres ...Scopemethod) error
	// UpdateByID 根据ID更新数据
	UpdateByID(id uint, m *T) error
	// Delete 删除数据
	Delete(wheres ...Scopemethod) error
	// DeleteByID 根据ID删除数据
	DeleteByID(id uint) error
	// ForcedDelete 强制删除数据
	ForcedDelete(wheres ...Scopemethod) error
	// ForcedDeleteByID 根据ID强制删除数据
	ForcedDeleteByID(id uint) error
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

	// OpenTrace 开启trace
	OpenTrace() IAction[T]
	// CloseTrace 关闭trace
	CloseTrace() IAction[T]
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

	// 开启trace
	enableTrace bool
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

// OpenTrace 开启trace
func OpenTrace[T any]() ActionOption[T] {
	return func(a *Action[T]) {
		a.enableTrace = true
	}
}

// WithDB 设置DB
func WithDB[T any](db *gorm.DB) ActionOption[T] {
	return func(a *Action[T]) {
		a.db = db
	}
}

// WithContext 设置Ctx
func WithContext[T any](ctx context.Context) ActionOption[T] {
	return func(a *Action[T]) {
		a.ctx = ctx
	}
}

// WithTabler 设置Tabler
func WithTabler[T any](tabler schema.Tabler) ActionOption[T] {
	return func(a *Action[T]) {
		a.Tabler = tabler
	}
}

func (a *Action[T]) OpenTrace() IAction[T] {
	a.enableTrace = true
	return a
}

func (a *Action[T]) CloseTrace() IAction[T] {
	a.enableTrace = false
	return a
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
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "First")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := a.DB().WithContext(ctx).Scopes(wheres...).First(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// FirstWithTrashed 查询单条数据(包含软删除数据)
func (a *Action[T]) FirstWithTrashed(wheres ...Scopemethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "FirstWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).First(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// FirstByID 根据ID查询单条数据
func (a *Action[T]) FirstByID(id uint) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "FirstByID")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := a.DB().WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// FirstByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
func (a *Action[T]) FirstByIDWithTrashed(id uint) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "FirstByIDWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := a.DB().WithContext(ctx).Unscoped().First(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// Last 查询单条数据
func (a *Action[T]) Last(wheres ...Scopemethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Last")
		defer span.End()
		ctx = _ctx
	}
	var m T

	if err := a.DB().WithContext(ctx).Scopes(wheres...).Last(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// LastWithTrashed 查询单条数据(包含软删除数据)
func (a *Action[T]) LastWithTrashed(wheres ...Scopemethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "LastWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var m T

	if err := a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).Last(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// LastByID 根据ID查询单条数据
func (a *Action[T]) LastByID(id uint) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "LastByID")
		defer span.End()
		ctx = _ctx
	}
	var m T

	if err := a.DB().WithContext(ctx).Last(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// LastByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
func (a *Action[T]) LastByIDWithTrashed(id uint) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "LastByIDWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var m T

	if err := a.DB().WithContext(ctx).Unscoped().Last(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// List 查询多条数据
func (a *Action[T]) List(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "List")
		defer span.End()
		ctx = _ctx
	}
	var ms []*T

	db := a.DB().WithContext(ctx).Scopes(wheres...)
	if pgInfo != nil {
		var total int64
		if err := db.WithContext(ctx).Count(&total).Error; err != nil {
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

// ListWithTrashed 查询多条数据(包含软删除数据)
func (a *Action[T]) ListWithTrashed(pgInfo Pagination, wheres ...Scopemethod) ([]*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "ListWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var ms []*T

	db := a.DB().WithContext(ctx).Unscoped().Scopes(wheres...)
	if pgInfo != nil {
		var total int64
		if err := db.WithContext(ctx).Count(&total).Error; err != nil {
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
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Count")
		defer span.End()
		ctx = _ctx
	}
	var total int64

	if err := a.DB().WithContext(ctx).Scopes(wheres...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// CountWithTrashed 查询数量(包含软删除数据)
func (a *Action[T]) CountWithTrashed(wheres ...Scopemethod) (int64, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "CountWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var total int64

	if err := a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Create 创建数据
func (a *Action[T]) Create(newModel *T) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Create")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Create(newModel).Error
}

// Update 更新数据
func (a *Action[T]) Update(newModel *T, wheres ...Scopemethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Update")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Scopes(wheres...).Updates(newModel).Error
}

// UpdateByID 根据ID更新数据
func (a *Action[T]) UpdateByID(id uint, newModel *T) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "UpdateByID")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Where("id = ?", id).Updates(newModel).Error
}

// Delete 删除数据
func (a *Action[T]) Delete(wheres ...Scopemethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Delete")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return a.DB().WithContext(ctx).Scopes(wheres...).Delete(&m).Error
}

// DeleteByID 根据ID删除数据
func (a *Action[T]) DeleteByID(id uint) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "DeleteByID")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return a.DB().WithContext(ctx).Delete(&m, id).Error
}

// ForcedDelete 强制删除数据
func (a *Action[T]) ForcedDelete(wheres ...Scopemethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "ForcedDelete")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).Delete(&m).Error
}

// ForcedDeleteByID 根据ID强制删除数据
func (a *Action[T]) ForcedDeleteByID(id uint) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "ForcedDeleteByID")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return a.DB().WithContext(ctx).Unscoped().Delete(&m, id).Error
}
