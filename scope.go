package query

import "gorm.io/gorm"

type ScopeMethod = func(db *gorm.DB) *gorm.DB

// WhereInColumn 通过字段名和值列表进行查询
func WhereInColumn[T any](column string, values ...T) ScopeMethod {
	return func(db *gorm.DB) *gorm.DB {
		idsLen := len(values)
		switch idsLen {
		case 0:
			return db
		default:
			return db.Where(column+" in (?)", values)
		}
	}
}

// WhereID 通过ID列表进行查询
func WhereID(ids ...uint) ScopeMethod {
	return WhereInColumn("id", ids...)
}

// WhereLikeKeyword 模糊查询
func WhereLikeKeyword(keyword string, columns ...string) ScopeMethod {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(columns) == 0 {
			return db
		}

		dbTmp := db
		for _, column := range columns {
			dbTmp = dbTmp.Or("`"+column+"` LIKE ?", keyword)
		}
		return db.Where(dbTmp)
	}
}

// BetweenColumn 通过字段名和值列表进行查询
func BetweenColumn[T any](column string, min, max T) ScopeMethod {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+" between ? and ?", min, max)
	}
}

// WhereColumn 通过字段名和值进行查询
func WhereColumn[T any](column string, val ...T) ScopeMethod {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column, val)
	}
}

// Paginate 分页
func Paginate(pgInfo Pagination) ScopeMethod {
	return func(db *gorm.DB) *gorm.DB {
		if pgInfo == nil {
			return db
		}
		return db.Limit(pgInfo.GetSize()).Offset((pgInfo.GetCurr() - 1) * pgInfo.GetSize())
	}
}

// WithTrashed 包含软删除数据
func WithTrashed(db *gorm.DB) *gorm.DB {
	return db.Unscoped()
}
