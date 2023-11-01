package query

import (
	"context"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type (
	// IAction 操作接口
	IAction[T any] interface {
		IOperation[T]
		IOperationX[T]
		IBind[T]

		DB() *gorm.DB
		Association() IAssociation
	}

	action[T any] struct {
		db    *gorm.DB
		ctx   context.Context
		table schema.Tabler

		association IAssociation

		// 开启trace
		enableTrace bool

		err error
	}

	ActionOption[T any] func(a *action[T])
)

func (a *action[T]) Err() error {
	return a.err
}

func (a *action[T]) FirstX(wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).First(&m).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) FirstWithTrashedX(wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).First(&m).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) FirstByIDX(id uint, wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).First(&m, id).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) FirstByIDWithTrashedX(id uint, wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).First(&m, id).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) LastX(wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Last(&m).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) LastWithTrashedX(wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).Last(&m).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) LastByIDX(id uint, wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Last(&m, id).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) LastByIDWithTrashedX(id uint, wheres ...ScopeMethod) *T {
	var m T
	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).Last(&m, id).Error; err != nil {
		a.err = err
		return nil
	}

	return &m
}

func (a *action[T]) ListX(pgInfo Pagination, wheres ...ScopeMethod) []*T {
	var ms []*T

	db := a.DB().WithContext(a.ctx).Scopes(wheres...)
	if pgInfo != nil {
		var total int64
		if err := db.WithContext(a.ctx).Count(&total).Error; err != nil {
			a.err = err
			return nil
		}
		pgInfo.SetTotal(total)
		db = db.Scopes(Paginate(pgInfo))
	}

	if err := db.Find(&ms).Error; err != nil {
		a.err = err
		return nil
	}

	return ms
}

func (a *action[T]) ListWithTrashedX(pgInfo Pagination, wheres ...ScopeMethod) []*T {
	var ms []*T

	db := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...)
	if pgInfo != nil {
		var total int64
		if err := db.WithContext(a.ctx).Count(&total).Error; err != nil {
			a.err = err
			return nil
		}
		pgInfo.SetTotal(total)
		db = db.Scopes(Paginate(pgInfo))
	}

	if err := db.Find(&ms).Error; err != nil {
		a.err = err
		return nil
	}

	return ms
}

func (a *action[T]) CountX(wheres ...ScopeMethod) int64 {
	var total int64

	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Count(&total).Error; err != nil {
		a.err = err
		return 0
	}

	return total
}

func (a *action[T]) CountWithTrashedX(wheres ...ScopeMethod) int64 {
	var total int64

	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).Count(&total).Error; err != nil {
		a.err = err
		return 0
	}

	return total
}

func (a *action[T]) CreateX(newModel *T) {
	if err := a.DB().WithContext(a.ctx).Create(newModel).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) UpdateX(newModel *T, wheres ...ScopeMethod) {
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Updates(newModel).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) UpdateMapX(newModel map[string]any, wheres ...ScopeMethod) {
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Updates(newModel).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) UpdateByIDX(id uint, newModel *T, wheres ...ScopeMethod) {
	if err := a.DB().WithContext(a.ctx).Scopes(append(wheres, WhereID(id))...).Updates(newModel).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) UpdateMapByIDX(id uint, newModel map[string]any, wheres ...ScopeMethod) {
	if err := a.DB().WithContext(a.ctx).Scopes(append(wheres, WhereID(id))...).Updates(newModel).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) DeleteX(wheres ...ScopeMethod) {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Delete(&m).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) DeleteByIDX(id uint, wheres ...ScopeMethod) {
	var m T
	if err := a.DB().WithContext(a.ctx).Scopes(wheres...).Delete(&m, id).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) ForcedDeleteX(wheres ...ScopeMethod) {
	var m T
	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).Delete(&m).Error; err != nil {
		a.err = err
	}
}

func (a *action[T]) ForcedDeleteByIDX(id uint, wheres ...ScopeMethod) {
	var m T
	if err := a.DB().WithContext(a.ctx).Unscoped().Scopes(wheres...).Delete(&m, id).Error; err != nil {
		a.err = err
	}
}

// NewAction 创建GORM操作接口实例
func NewAction[T any](opts ...ActionOption[T]) IAction[T] {
	ac := action[T]{
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(&ac)
	}

	if ac.table != nil {
		ac.db = ac.db.Table(ac.table.TableName())
	}

	ac.association = NewDefaultAssociation(ac.db)

	return &ac
}

// Association 转移到关联操作
func (a *action[T]) Association() IAssociation {
	return a.association
}

// OpenTrace 开启trace
func (a *action[T]) OpenTrace() IAction[T] {
	a.enableTrace = true
	return a
}

// CloseTrace 关闭trace
func (a *action[T]) CloseTrace() IAction[T] {
	a.enableTrace = false
	return a
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

// First 查询单条数据
func (a *action[T]) First(wheres ...ScopeMethod) (*T, error) {
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
func (a *action[T]) FirstWithTrashed(wheres ...ScopeMethod) (*T, error) {
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
func (a *action[T]) FirstByID(id uint, wheres ...ScopeMethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "FirstByID")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := a.DB().WithContext(ctx).Scopes(wheres...).First(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// FirstByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
func (a *action[T]) FirstByIDWithTrashed(id uint, wheres ...ScopeMethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "FirstByIDWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).First(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// Last 查询单条数据
func (a *action[T]) Last(wheres ...ScopeMethod) (*T, error) {
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
func (a *action[T]) LastWithTrashed(wheres ...ScopeMethod) (*T, error) {
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
func (a *action[T]) LastByID(id uint, wheres ...ScopeMethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "LastByID")
		defer span.End()
		ctx = _ctx
	}
	var m T

	if err := a.DB().WithContext(ctx).Scopes(wheres...).Last(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// LastByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
func (a *action[T]) LastByIDWithTrashed(id uint, wheres ...ScopeMethod) (*T, error) {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "LastByIDWithTrashed")
		defer span.End()
		ctx = _ctx
	}
	var m T

	if err := a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).Last(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

// List 查询多条数据
func (a *action[T]) List(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error) {
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
func (a *action[T]) ListWithTrashed(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error) {
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
func (a *action[T]) Count(wheres ...ScopeMethod) (int64, error) {
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
func (a *action[T]) CountWithTrashed(wheres ...ScopeMethod) (int64, error) {
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
func (a *action[T]) Create(newModel *T) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Create")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Create(newModel).Error
}

// Update 更新数据
func (a *action[T]) Update(newModel *T, wheres ...ScopeMethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Update")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Scopes(wheres...).Updates(newModel).Error
}

// UpdateMap 更新数据
func (a *action[T]) UpdateMap(newModel map[string]any, wheres ...ScopeMethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "Update")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Scopes(wheres...).Updates(newModel).Error
}

// UpdateByID 根据ID更新数据
func (a *action[T]) UpdateByID(id uint, newModel *T, wheres ...ScopeMethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "UpdateByID")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Scopes(append(wheres, WhereID(id))...).Updates(newModel).Error
}

// UpdateMapByID 根据ID更新数据
func (a *action[T]) UpdateMapByID(id uint, newModel map[string]any, wheres ...ScopeMethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "UpdateByID")
		defer span.End()
		ctx = _ctx
	}
	return a.DB().WithContext(ctx).Scopes(append(wheres, WhereID(id))...).Updates(newModel).Error
}

// Delete 删除数据
func (a *action[T]) Delete(wheres ...ScopeMethod) error {
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
func (a *action[T]) DeleteByID(id uint, wheres ...ScopeMethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "DeleteByID")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return a.DB().WithContext(ctx).Scopes(wheres...).Delete(&m, id).Error
}

// ForcedDelete 强制删除数据
func (a *action[T]) ForcedDelete(wheres ...ScopeMethod) error {
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
func (a *action[T]) ForcedDeleteByID(id uint, wheres ...ScopeMethod) error {
	ctx := a.ctx
	if a.enableTrace {
		_ctx, span := otel.Tracer("gorm-normalize").Start(a.ctx, "ForcedDeleteByID")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return a.DB().WithContext(ctx).Unscoped().Scopes(wheres...).Delete(&m, id).Error
}
