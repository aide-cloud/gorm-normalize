package query

import (
	"go.opentelemetry.io/otel"
)

type (
	operationMutationX[T any] struct {
		IBind[T]
		Tracer
		ICtx

		err error
	}

	OperationMutationXOption[T any] func(*operationMutationX[T])
)

func (l *operationMutationX[T]) CreateX(m *T) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "CreateX")
		defer span.End()
		ctx = _ctx
	}
	l.setErr(l.DB().WithContext(ctx).Create(m).Error)
}

func (l *operationMutationX[T]) BatchCreateX(m []*T, batchSize int) {
	if len(m) == 0 {
		return
	}

	if batchSize <= 0 {
		batchSize = 1000
	}

	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "BatchCreateX")
		defer span.End()
		ctx = _ctx
	}
	l.setErr(l.DB().WithContext(ctx).CreateInBatches(m, batchSize).Error)
}

func (l *operationMutationX[T]) UpdateX(m *T, wheres ...ScopeMethod) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "UpdateX")
		defer span.End()
		ctx = _ctx
	}
	l.setErr(l.DB().WithContext(ctx).Scopes(wheres...).Updates(m).Error)
}

func (l *operationMutationX[T]) UpdateMapX(m map[string]any, wheres ...ScopeMethod) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "UpdateMapX")
		defer span.End()
		ctx = _ctx
	}
	l.setErr(l.DB().WithContext(ctx).Scopes(wheres...).Updates(m).Error)
}

func (l *operationMutationX[T]) UpdateByIDX(id uint, m *T, wheres ...ScopeMethod) {
	l.UpdateX(m, append(wheres, WhereID(id))...)
}

func (l *operationMutationX[T]) UpdateMapByIDX(id uint, m map[string]any, wheres ...ScopeMethod) {
	l.UpdateMapX(m, append(wheres, WhereID(id))...)
}

func (l *operationMutationX[T]) DeleteX(wheres ...ScopeMethod) {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "DeleteX")
		defer span.End()
		ctx = _ctx
	}
	var m T
	l.setErr(l.DB().WithContext(ctx).Scopes(wheres...).Delete(&m).Error)
}

func (l *operationMutationX[T]) DeleteByIDX(id uint, wheres ...ScopeMethod) {
	l.DeleteX(append(wheres, WhereID(id))...)
}

func (l *operationMutationX[T]) ForcedDeleteX(wheres ...ScopeMethod) {
	l.DeleteX(append(wheres, WithTrashed)...)
}

func (l *operationMutationX[T]) ForcedDeleteByIDX(id uint, wheres ...ScopeMethod) {
	l.DeleteX(append(wheres, WhereID(id), WithTrashed)...)
}

func (l *operationMutationX[T]) setErr(err error) {
	if err != nil {
		l.err = err
	}
}

func (l *operationMutationX[T]) GetMutationErr() error {
	return l.err
}

func defaultOperationMutationX[T any]() *operationMutationX[T] {
	return &operationMutationX[T]{}
}

// NewOperationMutationX 实例化操作
func NewOperationMutationX[T any](opts ...OperationMutationXOption[T]) IOperationMutationX[T] {
	o := defaultOperationMutationX[T]()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithOperationMutationXBind 设置绑定
func WithOperationMutationXBind[T any](bind IBind[T]) OperationMutationXOption[T] {
	return func(o *operationMutationX[T]) {
		o.IBind = bind
	}
}

// WithOperationMutationXTracer 设置跟踪
func WithOperationMutationXTracer[T any](tracer Tracer) OperationMutationXOption[T] {
	return func(o *operationMutationX[T]) {
		o.Tracer = tracer
	}
}

// WithOperationMutationXCtx 设置上下文
func WithOperationMutationXCtx[T any](ctx ICtx) OperationMutationXOption[T] {
	return func(o *operationMutationX[T]) {
		o.ICtx = ctx
	}
}
