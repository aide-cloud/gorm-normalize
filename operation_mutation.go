package query

import (
	"go.opentelemetry.io/otel"
)

type (
	operationMutation[T any] struct {
		IBind[T]
		Tracer
		ICtx
	}

	OperationMutationOption[T any] func(*operationMutation[T])
)

func (l *operationMutation[T]) Create(m *T) error {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "Create")
		defer span.End()
		ctx = _ctx
	}

	return l.DB().WithContext(ctx).Create(m).Error
}

func (l *operationMutation[T]) BatchCreate(m []*T, batchSize int) error {
	if len(m) == 0 {
		return nil
	}

	if batchSize <= 0 {
		batchSize = 1000
	}

	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "BatchCreate")
		defer span.End()
		ctx = _ctx
	}
	return l.DB().WithContext(ctx).CreateInBatches(m, batchSize).Error
}

func (l *operationMutation[T]) Update(m *T, wheres ...ScopeMethod) error {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "Update")
		defer span.End()
		ctx = _ctx
	}
	return l.DB().WithContext(ctx).Scopes(wheres...).Updates(m).Error
}

func (l *operationMutation[T]) UpdateMap(m map[string]any, wheres ...ScopeMethod) error {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "UpdateMap")
		defer span.End()
		ctx = _ctx
	}
	return l.DB().WithContext(ctx).Scopes(wheres...).Updates(m).Error
}

func (l *operationMutation[T]) UpdateByID(id uint32, m *T, wheres ...ScopeMethod) error {
	return l.Update(m, append(wheres, WhereID(id))...)
}

func (l *operationMutation[T]) UpdateMapByID(id uint32, m map[string]any, wheres ...ScopeMethod) error {
	return l.UpdateMap(m, append(wheres, WhereID(id))...)
}

func (l *operationMutation[T]) Delete(wheres ...ScopeMethod) error {
	ctx := l.GetCtx()
	if l.IsEnableTrace() {
		_ctx, span := otel.Tracer("gorm-normalize").Start(ctx, "Delete")
		defer span.End()
		ctx = _ctx
	}
	var m T
	return l.DB().WithContext(ctx).Scopes(wheres...).Delete(&m).Error
}

func (l *operationMutation[T]) DeleteByID(id uint32, wheres ...ScopeMethod) error {
	return l.Delete(append(wheres, WhereID(id))...)
}

func (l *operationMutation[T]) ForcedDelete(wheres ...ScopeMethod) error {
	return l.Delete(append(wheres, WithTrashed)...)
}

func (l *operationMutation[T]) ForcedDeleteByID(id uint32, wheres ...ScopeMethod) error {
	return l.Delete(append(wheres, WhereID(id), WithTrashed)...)
}

func defaultOperationMutation[T any]() *operationMutation[T] {
	return &operationMutation[T]{}
}

// NewOperationMutation 实例化操作
func NewOperationMutation[T any](opts ...OperationMutationOption[T]) IOperationMutation[T] {
	o := defaultOperationMutation[T]()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithOperationMutationTracer 设置跟踪
func WithOperationMutationTracer[T any](tracer Tracer) OperationMutationOption[T] {
	return func(o *operationMutation[T]) {
		o.Tracer = tracer
	}
}

// WithOperationMutationICtx 设置上下文
func WithOperationMutationICtx[T any](ctx ICtx) OperationMutationOption[T] {
	return func(o *operationMutation[T]) {
		o.ICtx = ctx
	}
}

// WithOperationMutationIBind 设置数据库
func WithOperationMutationIBind[T any](bind IBind[T]) OperationMutationOption[T] {
	return func(o *operationMutation[T]) {
		o.IBind = bind
	}
}
