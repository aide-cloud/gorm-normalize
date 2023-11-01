package query

import (
	"context"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var _ IAction[any] = (*action[any])(nil)

type IOperation[T any] interface {
	// First 查询单条数据
	First(wheres ...ScopeMethod) (*T, error)
	// FirstWithTrashed 查询单条数据(包含软删除数据)
	FirstWithTrashed(wheres ...ScopeMethod) (*T, error)
	// FirstByID 根据ID查询单条数据
	FirstByID(id uint, wheres ...ScopeMethod) (*T, error)
	// FirstByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
	FirstByIDWithTrashed(id uint, wheres ...ScopeMethod) (*T, error)
	// Last 查询单条数据
	Last(wheres ...ScopeMethod) (*T, error)
	// LastWithTrashed 查询单条数据(包含软删除数据)
	LastWithTrashed(wheres ...ScopeMethod) (*T, error)
	// LastByID 根据ID查询单条数据
	LastByID(id uint, wheres ...ScopeMethod) (*T, error)
	// LastByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
	LastByIDWithTrashed(id uint, wheres ...ScopeMethod) (*T, error)
	// List 查询多条数据
	List(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error)
	// ListWithTrashed 查询多条数据(包含软删除数据)
	ListWithTrashed(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error)
	// Count 查询数量
	Count(wheres ...ScopeMethod) (int64, error)
	// CountWithTrashed 查询数量(包含软删除数据)
	CountWithTrashed(wheres ...ScopeMethod) (int64, error)
	// Create 创建数据
	Create(m *T) error
	// Update 更新数据
	Update(m *T, wheres ...ScopeMethod) error
	// UpdateMap 通过map更新数据
	UpdateMap(m map[string]any, wheres ...ScopeMethod) error
	// UpdateByID 根据ID更新数据
	UpdateByID(id uint, m *T, wheres ...ScopeMethod) error
	// UpdateMapByID 根据ID更新数据
	UpdateMapByID(id uint, m map[string]any, wheres ...ScopeMethod) error
	// Delete 删除数据
	Delete(wheres ...ScopeMethod) error
	// DeleteByID 根据ID删除数据
	DeleteByID(id uint, wheres ...ScopeMethod) error
	// ForcedDelete 强制删除数据
	ForcedDelete(wheres ...ScopeMethod) error
	// ForcedDeleteByID 根据ID强制删除数据
	ForcedDeleteByID(id uint, wheres ...ScopeMethod) error
}

type AssociationKey string

// IAssociation 关联操作
type IAssociation interface {
	// Append 添加关联
	Append(associationKey AssociationKey, list ...schema.Tabler) error
	// Replace 替换关联
	Replace(associationKey AssociationKey, list ...schema.Tabler) error
	// Delete 删除关联
	Delete(associationKey AssociationKey, list ...schema.Tabler) error
	// Clear 清除关联
	Clear(associationKey AssociationKey) error
	// Count 关联数量
	Count(associationKey AssociationKey) int64
}

// defaultAssociation 默认关联操作实现
type defaultAssociation struct {
	db *gorm.DB
}

// Append 添加关联
func (l *defaultAssociation) Append(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Append(list)
}

// Replace 替换关联
func (l *defaultAssociation) Replace(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Replace(list)
}

// Delete 删除关联
func (l *defaultAssociation) Delete(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Delete(list)
}

// Clear 清除关联
func (l *defaultAssociation) Clear(associationKey AssociationKey) error {
	return l.db.Association(string(associationKey)).Clear()
}

// Count 统计关联数量
func (l *defaultAssociation) Count(associationKey AssociationKey) int64 {
	return l.db.Association(string(associationKey)).Count()
}

// NewDefaultAssociation 创建默认关联
func NewDefaultAssociation(db *gorm.DB) IAssociation {
	return &defaultAssociation{
		db: db,
	}
}

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

// IAction 操作接口
type IAction[T any] interface {
	IOperation[T]
	IBind[T]

	DB() *gorm.DB
	Association() IAssociation
}

type action[T any] struct {
	db  *gorm.DB
	ctx context.Context
	schema.Tabler

	association IAssociation

	// 开启trace
	enableTrace bool
}

type ActionOption[T any] func(a *action[T])

// NewAction 创建GORM操作接口实例
func NewAction[T any](opts ...ActionOption[T]) IAction[T] {
	ac := action[T]{
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(&ac)
	}

	if ac.Tabler != nil {
		ac.db = ac.db.Table(ac.Tabler.TableName())
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
	if a.Tabler != nil {
		return a.db.Table(a.Tabler.TableName())
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
	a.Tabler = tabler
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
