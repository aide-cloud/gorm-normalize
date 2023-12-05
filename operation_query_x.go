package query

import (
	"go.opentelemetry.io/otel"
)

type (
	operationQueryX[T any] struct {
		IBind[T]
		Tracer
		ICtx
		err error
	}

	OperationQueryXOption[T any] func(*operationQueryX[T])
)

// setErr 设置错误
func (l *operationQueryX[T]) setErr(err error) {
	if err != nil {
		l.err = err
	}
}

func (l *operationQueryX[T]) FirstX(wheres ...ScopeMethod) *T {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "FirstX")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := l.DB().WithContext(ctx).Scopes(wheres...).First(&m).Error; err != nil {
		l.setErr(err)
		return nil
	}
	return &m
}

func (l *operationQueryX[T]) FirstWithTrashedX(wheres ...ScopeMethod) *T {
	return l.FirstX(append(wheres, WithTrashed)...)
}

func (l *operationQueryX[T]) FirstByIDX(id uint32, wheres ...ScopeMethod) *T {
	return l.FirstX(append(wheres, WhereID(id))...)
}

func (l *operationQueryX[T]) FirstByIDWithTrashedX(id uint32, wheres ...ScopeMethod) *T {
	return l.FirstX(append(wheres, WhereID(id), WithTrashed)...)
}

func (l *operationQueryX[T]) LastX(wheres ...ScopeMethod) *T {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "LastX")
		defer span.End()
		ctx = _ctx
	}
	var m T
	if err := l.DB().WithContext(ctx).Scopes(wheres...).Last(&m).Error; err != nil {
		l.setErr(err)
		return nil
	}
	return &m
}

func (l *operationQueryX[T]) LastWithTrashedX(wheres ...ScopeMethod) *T {
	return l.LastX(append(wheres, WithTrashed)...)
}

func (l *operationQueryX[T]) LastByIDX(id uint32, wheres ...ScopeMethod) *T {
	return l.LastX(append(wheres, WhereID(id))...)
}

func (l *operationQueryX[T]) LastByIDWithTrashedX(id uint32, wheres ...ScopeMethod) *T {
	return l.LastX(append(wheres, WhereID(id), WithTrashed)...)
}

func (l *operationQueryX[T]) ListX(pgInfo Pagination, wheres ...ScopeMethod) []*T {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "ListX")
		defer span.End()
		ctx = _ctx
	}
	var ms []*T

	db := l.DB().WithContext(ctx).Scopes(wheres...)
	if pgInfo != nil {
		pgInfo.SetTotal(l.CountX(wheres...))
		db = db.Scopes(Paginate(pgInfo))
	}

	if err := db.Find(&ms).Error; err != nil {
		l.setErr(err)
		return nil
	}

	return ms
}

func (l *operationQueryX[T]) ListWithTrashedX(pgInfo Pagination, wheres ...ScopeMethod) []*T {
	return l.ListX(pgInfo, append(wheres, WithTrashed)...)
}

func (l *operationQueryX[T]) CountX(wheres ...ScopeMethod) int64 {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "CountX")
		defer span.End()
		ctx = _ctx
	}
	var total int64
	if err := l.DB().WithContext(ctx).Scopes(wheres...).Count(&total).Error; err != nil {
		l.setErr(err)
		return 0
	}
	return total
}

func (l *operationQueryX[T]) CountWithTrashedX(wheres ...ScopeMethod) int64 {
	return l.CountX(append(wheres, WithTrashed)...)
}

func (l *operationQueryX[T]) GetQueryErr() error {
	return l.err
}

func defaultOperationQueryX[T any]() *operationQueryX[T] {
	return &operationQueryX[T]{}
}

// NewOperationQueryX 实例化查询操作
func NewOperationQueryX[T any](opts ...OperationQueryXOption[T]) IOperationQueryX[T] {
	o := defaultOperationQueryX[T]()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithOperationQueryXIBind 设置IBind
func WithOperationQueryXIBind[T any](b IBind[T]) OperationQueryXOption[T] {
	return func(o *operationQueryX[T]) {
		o.IBind = b
	}
}

// WithOperationQueryXTracer 设置Tracer
func WithOperationQueryXTracer[T any](t Tracer) OperationQueryXOption[T] {
	return func(o *operationQueryX[T]) {
		o.Tracer = t
	}
}

// WithOperationQueryXICtx 设置ICtx
func WithOperationQueryXICtx[T any](c ICtx) OperationQueryXOption[T] {
	return func(o *operationQueryX[T]) {
		o.ICtx = c
	}
}
