package query

// IOperationQueryX 扩展查询操作, 不返回error
type IOperationQueryX[T any] interface {
	// FirstX 查询单条数据
	FirstX(wheres ...ScopeMethod) *T
	// FirstWithTrashedX 查询单条数据(包含软删除数据)
	FirstWithTrashedX(wheres ...ScopeMethod) *T
	// FirstByIDX 根据ID查询单条数据
	FirstByIDX(id uint32, wheres ...ScopeMethod) *T
	// FirstByIDWithTrashedX 根据ID查询单条数据(包含软删除数据)
	FirstByIDWithTrashedX(id uint32, wheres ...ScopeMethod) *T
	// LastX 查询单条数据
	LastX(wheres ...ScopeMethod) *T
	// LastWithTrashedX 查询单条数据(包含软删除数据)
	LastWithTrashedX(wheres ...ScopeMethod) *T
	// LastByIDX 根据ID查询单条数据
	LastByIDX(id uint32, wheres ...ScopeMethod) *T
	// LastByIDWithTrashedX 根据ID查询单条数据(包含软删除数据)
	LastByIDWithTrashedX(id uint32, wheres ...ScopeMethod) *T
	// ListX 查询多条数据
	ListX(pgInfo Pagination, wheres ...ScopeMethod) []*T
	// ListWithTrashedX 查询多条数据(包含软删除数据)
	ListWithTrashedX(pgInfo Pagination, wheres ...ScopeMethod) []*T
	// CountX 查询数量
	CountX(wheres ...ScopeMethod) int64
	// CountWithTrashedX 查询数量(包含软删除数据)
	CountWithTrashedX(wheres ...ScopeMethod) int64

	GetQueryErr() error
}

// IOperationMutationX 变更操作, 不返回error
type IOperationMutationX[T any] interface {
	// CreateX 创建数据
	CreateX(m *T)
	// BatchCreateX 批量创建数据
	BatchCreateX(m []*T, batchSize int)
	// UpdateX 更新数据
	UpdateX(m *T, wheres ...ScopeMethod)
	// UpdateMapX 通过map更新数据
	UpdateMapX(m map[string]any, wheres ...ScopeMethod)
	// UpdateByIDX 根据ID更新数据
	UpdateByIDX(id uint32, m *T, wheres ...ScopeMethod)
	// UpdateMapByIDX 根据ID更新数据
	UpdateMapByIDX(id uint32, m map[string]any, wheres ...ScopeMethod)
	// DeleteX 删除数据
	DeleteX(wheres ...ScopeMethod)
	// DeleteByIDX 根据ID删除数据
	DeleteByIDX(id uint32, wheres ...ScopeMethod)
	// ForcedDeleteX 强制删除数据
	ForcedDeleteX(wheres ...ScopeMethod)
	// ForcedDeleteByIDX 根据ID强制删除数据
	ForcedDeleteByIDX(id uint32, wheres ...ScopeMethod)

	GetMutationErr() error
}

// IOperationX 扩展操作, 不返回error
type IOperationX[T any] interface {
	IOperationQueryX[T]
	IOperationMutationX[T]
	Err() error
}
