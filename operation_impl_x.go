package query

type (
	operationXImpl[T any] struct {
		IOperationQueryX[T]
		IOperationMutationX[T]
	}

	OperationXImplOption[T any] func(*operationXImpl[T])
)

// defaultOperationXImpl 默认操作实现
func defaultOperationXImpl[T any]() *operationXImpl[T] {
	return &operationXImpl[T]{
		IOperationQueryX:    NewOperationQueryX[T](),
		IOperationMutationX: NewOperationMutationX[T](),
	}
}

// NewOperationXImpl 实例化操作
func NewOperationXImpl[T any](opts ...OperationXImplOption[T]) IOperationX[T] {
	o := defaultOperationXImpl[T]()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Err 错误
func (o *operationXImpl[T]) Err() error {
	if o.IOperationQueryX.GetQueryErr() != nil {
		return o.IOperationQueryX.GetQueryErr()
	}
	return o.IOperationMutationX.GetMutationErr()
}

// WithOperationQueryX 设置查询
func WithOperationQueryX[T any](query IOperationQueryX[T]) OperationXImplOption[T] {
	return func(o *operationXImpl[T]) {
		o.IOperationQueryX = query
	}
}

// WithOperationMutationX 设置变更
func WithOperationMutationX[T any](mutation IOperationMutationX[T]) OperationXImplOption[T] {
	return func(o *operationXImpl[T]) {
		o.IOperationMutationX = mutation
	}
}
