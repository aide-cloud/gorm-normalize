package query

// IOperationQuery 扩展查询操作, 返回error
type IOperationQuery[T any] interface {
	// First 查询单条数据
	First(wheres ...ScopeMethod) (*T, error)
	// FirstWithTrashed 查询单条数据(包含软删除数据)
	FirstWithTrashed(wheres ...ScopeMethod) (*T, error)
	// FirstByID 根据ID查询单条数据
	FirstByID(id uint32, wheres ...ScopeMethod) (*T, error)
	// FirstByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
	FirstByIDWithTrashed(id uint32, wheres ...ScopeMethod) (*T, error)
	// Last 查询单条数据
	Last(wheres ...ScopeMethod) (*T, error)
	// LastWithTrashed 查询单条数据(包含软删除数据)
	LastWithTrashed(wheres ...ScopeMethod) (*T, error)
	// LastByID 根据ID查询单条数据
	LastByID(id uint32, wheres ...ScopeMethod) (*T, error)
	// LastByIDWithTrashed 根据ID查询单条数据(包含软删除数据)
	LastByIDWithTrashed(id uint32, wheres ...ScopeMethod) (*T, error)
	// List 查询多条数据
	List(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error)
	// ListWithTrashed 查询多条数据(包含软删除数据)
	ListWithTrashed(pgInfo Pagination, wheres ...ScopeMethod) ([]*T, error)
	// Count 查询数量
	Count(wheres ...ScopeMethod) (int64, error)
	// CountWithTrashed 查询数量(包含软删除数据)
	CountWithTrashed(wheres ...ScopeMethod) (int64, error)
}

// IOperationMutation 变更操作, 返回error
type IOperationMutation[T any] interface {
	// Create 创建数据
	Create(m *T) error
	// BatchCreate 批量创建数据
	BatchCreate(m []*T, max int) error
	// Update 更新数据
	Update(m *T, wheres ...ScopeMethod) error
	// UpdateMap 通过map更新数据
	UpdateMap(m map[string]any, wheres ...ScopeMethod) error
	// UpdateByID 根据ID更新数据
	UpdateByID(id uint32, m *T, wheres ...ScopeMethod) error
	// UpdateMapByID 根据ID更新数据
	UpdateMapByID(id uint32, m map[string]any, wheres ...ScopeMethod) error
	// Delete 删除数据
	Delete(wheres ...ScopeMethod) error
	// DeleteByID 根据ID删除数据
	DeleteByID(id uint32, wheres ...ScopeMethod) error
	// ForcedDelete 强制删除数据
	ForcedDelete(wheres ...ScopeMethod) error
	// ForcedDeleteByID 根据ID强制删除数据
	ForcedDeleteByID(id uint32, wheres ...ScopeMethod) error
}

type IOperation[T any] interface {
	IOperationQuery[T]
	IOperationMutation[T]
}
