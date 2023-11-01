package query

// IOperationQuery 扩展查询操作, 返回error
type IOperationQuery[T any] interface {
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
}

// IMutation 变更操作, 返回error
type IMutation[T any] interface {
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

type IOperation[T any] interface {
	IOperationQuery[T]
	IMutation[T]
}

// IOperationQueryX 扩展查询操作, 不返回error
type IOperationQueryX[T any] interface {
	// FirstX 查询单条数据
	FirstX(wheres ...ScopeMethod) *T
	// FirstWithTrashedX 查询单条数据(包含软删除数据)
	FirstWithTrashedX(wheres ...ScopeMethod) *T
	// FirstByIDX 根据ID查询单条数据
	FirstByIDX(id uint, wheres ...ScopeMethod) *T
	// FirstByIDWithTrashedX 根据ID查询单条数据(包含软删除数据)
	FirstByIDWithTrashedX(id uint, wheres ...ScopeMethod) *T
	// LastX 查询单条数据
	LastX(wheres ...ScopeMethod) *T
	// LastWithTrashedX 查询单条数据(包含软删除数据)
	LastWithTrashedX(wheres ...ScopeMethod) *T
	// LastByIDX 根据ID查询单条数据
	LastByIDX(id uint, wheres ...ScopeMethod) *T
	// LastByIDWithTrashedX 根据ID查询单条数据(包含软删除数据)
	LastByIDWithTrashedX(id uint, wheres ...ScopeMethod) *T
	// ListX 查询多条数据
	ListX(pgInfo Pagination, wheres ...ScopeMethod) []*T
	// ListWithTrashedX 查询多条数据(包含软删除数据)
	ListWithTrashedX(pgInfo Pagination, wheres ...ScopeMethod) []*T
	// CountX 查询数量
	CountX(wheres ...ScopeMethod) int64
	// CountWithTrashedX 查询数量(包含软删除数据)
	CountWithTrashedX(wheres ...ScopeMethod) int64
}

// IMutationX 变更操作, 不返回error
type IMutationX[T any] interface {
	// CreateX 创建数据
	CreateX(m *T)
	// UpdateX 更新数据
	UpdateX(m *T, wheres ...ScopeMethod)
	// UpdateMapX 通过map更新数据
	UpdateMapX(m map[string]any, wheres ...ScopeMethod)
	// UpdateByIDX 根据ID更新数据
	UpdateByIDX(id uint, m *T, wheres ...ScopeMethod)
	// UpdateMapByIDX 根据ID更新数据
	UpdateMapByIDX(id uint, m map[string]any, wheres ...ScopeMethod)
	// DeleteX 删除数据
	DeleteX(wheres ...ScopeMethod)
	// DeleteByIDX 根据ID删除数据
	DeleteByIDX(id uint, wheres ...ScopeMethod)
	// ForcedDeleteX 强制删除数据
	ForcedDeleteX(wheres ...ScopeMethod)
	// ForcedDeleteByIDX 根据ID强制删除数据
	ForcedDeleteByIDX(id uint, wheres ...ScopeMethod)
}

// IOperationX 扩展操作, 不返回error
type IOperationX[T any] interface {
	IOperationQueryX[T]
	IMutationX[T]
	Err() error
}
