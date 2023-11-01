package query

type operationQuery[T any] struct {
}

func (l *operationQuery[T]) First(wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) FirstWithTrashed(wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) FirstByID(id uint, wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) FirstByIDWithTrashed(id uint, wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) Last(wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) LastWithTrashed(wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) LastByID(id uint, wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) LastByIDWithTrashed(id uint, wheres ...ScopeMethod) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) List(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) ListWithTrashed(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) Count(wheres ...ScopeMethod) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (l *operationQuery[T]) CountWithTrashed(wheres ...ScopeMethod) (int64, error) {
	//TODO implement me
	panic("implement me")
}

// NewOperationQuery 实例化查询操作
func NewOperationQuery[T any]() IOperationQuery[T] {
	return &operationQuery[T]{}
}
