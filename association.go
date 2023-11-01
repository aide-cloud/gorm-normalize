package query

import (
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type AssociationKey string

// IAssociation 关联操作
type IAssociation interface {
	// Append 添加关联
	Append(associationKey AssociationKey, list ...schema.Tabler) error
	// Replace 替换关联
	Replace(associationKey AssociationKey, list ...schema.Tabler) error
	// Delete 删除关联
	Delete(associationKey AssociationKey, list ...schema.Tabler) error
	// Clear 清除关联
	Clear(associationKey AssociationKey) error
	// Count 关联数量
	Count(associationKey AssociationKey) int64
}

// defaultAssociation 默认关联操作实现
type defaultAssociation struct {
	db *gorm.DB
}

// Append 添加关联
func (l *defaultAssociation) Append(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Append(list)
}

// Replace 替换关联
func (l *defaultAssociation) Replace(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Replace(list)
}

// Delete 删除关联
func (l *defaultAssociation) Delete(associationKey AssociationKey, list ...schema.Tabler) error {
	return l.db.Association(string(associationKey)).Delete(list)
}

// Clear 清除关联
func (l *defaultAssociation) Clear(associationKey AssociationKey) error {
	return l.db.Association(string(associationKey)).Clear()
}

// Count 统计关联数量
func (l *defaultAssociation) Count(associationKey AssociationKey) int64 {
	return l.db.Association(string(associationKey)).Count()
}

// NewDefaultAssociation 创建默认关联
func NewDefaultAssociation(db *gorm.DB) IAssociation {
	return &defaultAssociation{
		db: db,
	}
}
