package query

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type AssociationKey string

// IAssociation 关联操作
type IAssociation interface {
	// AssociationAppend 添加关联
	AssociationAppend(associationKey AssociationKey, list ...schema.Tabler) error
	// AssociationReplace 替换关联
	AssociationReplace(associationKey AssociationKey, list ...schema.Tabler) error
	// AssociationDelete 删除关联
	AssociationDelete(associationKey AssociationKey, list ...schema.Tabler) error
	// AssociationClear 清除关联
	AssociationClear(associationKey AssociationKey) error
	// AssociationCount 关联数量
	AssociationCount(associationKey AssociationKey) int64
}

// defaultAssociation 默认关联操作实现
type defaultAssociation struct {
	db *gorm.DB
}

// AssociationAppend 添加关联
func (l *defaultAssociation) AssociationAppend(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Append(list)
}

// AssociationReplace 替换关联
func (l *defaultAssociation) AssociationReplace(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Replace(list)
}

// AssociationDelete 删除关联
func (l *defaultAssociation) AssociationDelete(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Delete(list)
}

// AssociationClear 清除关联
func (l *defaultAssociation) AssociationClear(associationKey AssociationKey) error {
	return l.db.Association(string(associationKey)).Clear()
}

// AssociationCount 统计关联数量
func (l *defaultAssociation) AssociationCount(associationKey AssociationKey) int64 {
	return l.db.Association(string(associationKey)).Count()
}

// NewDefaultAssociation 创建默认关联
func NewDefaultAssociation(db *gorm.DB) IAssociation {
	return &defaultAssociation{
		db: db,
	}
}
