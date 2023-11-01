package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// OpenTrace 开启trace
func OpenTrace[T any]() ActionOption[T] {
	return func(a *action[T]) {
		a.enableTrace = true
	}
}

// WithDB 设置DB
func WithDB[T any](db *gorm.DB) ActionOption[T] {
	return func(a *action[T]) {
		a.db = db
	}
}

// WithContext 设置Ctx
func WithContext[T any](ctx context.Context) ActionOption[T] {
	return func(a *action[T]) {
		a.ctx = ctx
	}
}

// WithTable 设置Table
func WithTable[T any](table schema.Tabler) ActionOption[T] {
	return func(a *action[T]) {
		a.table = table
	}
}
