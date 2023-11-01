package query

type (
	OperationImpl[T any] struct {
		IOperationQuery[T]
		IOperationMutation[T]
	}

	OperationImplOption[T any] func(*OperationImpl[T])
)

// defaultOperationImpl 默认操作实现
func defaultOperationImpl[T any]() *OperationImpl[T] {
	return &OperationImpl[T]{
		IOperationQuery:    NewOperationQuery[T](),
		IOperationMutation: NewOperationMutation[T](),
	}
}

// NewOperationImpl 实例化操作
func NewOperationImpl[T any](opts ...OperationImplOption[T]) IOperation[T] {
	o := defaultOperationImpl[T]()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// WithOperationQuery 设置查询
func WithOperationQuery[T any](query IOperationQuery[T]) OperationImplOption[T] {
	return func(o *OperationImpl[T]) {
		o.IOperationQuery = query
	}
}

// WithOperationMutation 设置变更
func WithOperationMutation[T any](mutation IOperationMutation[T]) OperationImplOption[T] {
	return func(o *OperationImpl[T]) {
		o.IOperationMutation = mutation
	}
}
