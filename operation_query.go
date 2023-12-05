package query

import (
	"go.opentelemetry.io/otel"
)

type (
	operationQuery[T any] struct {
		IBind[T]
		Tracer
		ICtx
	}

	OperationQueryOption[T any] func(*operationQuery[T])
)

func defaultOperationQuery[T any]() *operationQuery[T] {
	return &operationQuery[T]{}
}

// NewOperationQuery 实例化查询操作
func NewOperationQuery[T any](opts ...OperationQueryOption[T]) IOperationQuery[T] {
	o := defaultOperationQuery[T]()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithOperationQueryIBind 设置IBind
func WithOperationQueryIBind[T any](b IBind[T]) OperationQueryOption[T] {
	return func(o *operationQuery[T]) {
		o.IBind = b
	}
}

// WithOperationQueryTracer 设置Tracer
func WithOperationQueryTracer[T any](t Tracer) OperationQueryOption[T] {
	return func(o *operationQuery[T]) {
		o.Tracer = t
	}
}

// WithOperationQueryICtx 设置ICtx
func WithOperationQueryICtx[T any](c ICtx) OperationQueryOption[T] {
	return func(o *operationQuery[T]) {
		o.ICtx = c
	}
}

func (l *operationQuery[T]) First(wheres ...ScopeMethod) (*T, error) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "First")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := l.DB().WithContext(ctx).Scopes(wheres...).First(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

func (l *operationQuery[T]) FirstWithTrashed(wheres ...ScopeMethod) (*T, error) {
	return l.First(append(wheres, WithTrashed)...)
}

func (l *operationQuery[T]) FirstByID(id uint32, wheres ...ScopeMethod) (*T, error) {
	return l.First(append(wheres, WhereID(id))...)
}

func (l *operationQuery[T]) FirstByIDWithTrashed(id uint32, wheres ...ScopeMethod) (*T, error) {
	return l.First(append(wheres, WhereID(id), WithTrashed)...)
}

func (l *operationQuery[T]) Last(wheres ...ScopeMethod) (*T, error) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "Last")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := l.DB().WithContext(ctx).Scopes(wheres...).Last(&m).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

func (l *operationQuery[T]) LastWithTrashed(wheres ...ScopeMethod) (*T, error) {
	return l.Last(append(wheres, WithTrashed)...)
}

func (l *operationQuery[T]) LastByID(id uint32, wheres ...ScopeMethod) (*T, error) {
	return l.Last(append(wheres, WhereID(id))...)
}

func (l *operationQuery[T]) LastByIDWithTrashed(id uint32, wheres ...ScopeMethod) (*T, error) {
	return l.Last(append(wheres, WhereID(id), WithTrashed)...)
}

func (l *operationQuery[T]) List(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "List")
		defer span.End()
		ctx = _ctx
	}
	var ms []*T

	db := l.DB().WithContext(ctx).Scopes(wheres...)
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

func (l *operationQuery[T]) ListWithTrashed(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error) {
	return l.List(pgInfo, append(wheres, WithTrashed)...)
}

func (l *operationQuery[T]) Count(wheres ...ScopeMethod) (int64, error) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "Count")
		defer span.End()
		ctx = _ctx
	}
	var total int64
	if err := l.DB().WithContext(ctx).Scopes(wheres...).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (l *operationQuery[T]) CountWithTrashed(wheres ...ScopeMethod) (int64, error) {
	return l.Count(append(wheres, WithTrashed)...)
}
