package query

import "gorm.io/gorm"

type Scopemethod = func(db *gorm.DB) *gorm.DB

// WhereInColumn 通过字段名和值列表进行查询
func WhereInColumn[T any](column string, vals ...T) Scopemethod {
	return func(db *gorm.DB) *gorm.DB {
		idsLen := len(vals)
		switch idsLen {
		case 0:
			return db
		default:
			return db.Where(column+" in (?)", vals)
		}
	}
}

// WhereID 通过ID列表进行查询
func WhereID(ids ...uint) Scopemethod {
	return WhereInColumn("id", ids...)
}

// WhereLikeKeyword 模糊查询
func WhereLikeKeyword(keyword string, columns ...string) Scopemethod {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(columns) == 0 {
			return db
		}

		dbTmp := db
		likeKeyword := "%" + keyword + "%"
		for _, column := range columns {
			dbTmp = dbTmp.Or("`"+column+"` LIKE ?", likeKeyword)
		}
		return db.Where(dbTmp)
	}
}

// BetweenColumn 通过字段名和值列表进行查询
func BetweenColumn[T any](column string, min, max T) Scopemethod {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column+" between ? and ?", min, max)
	}
}

// WhereColumn 通过字段名和值进行查询
func WhereColumn[T any](column string, val ...T) Scopemethod {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(column, val)
	}
}

// Paginate 分页
func Paginate(pgInfo Pagination) Scopemethod {
	return func(db *gorm.DB) *gorm.DB {
		if pgInfo == nil {
			return db
		}
		return db.Limit(pgInfo.Size()).Offset((pgInfo.Page() - 1) * pgInfo.Size())
	}
}
