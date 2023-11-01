package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// WithTracer 设置Tracer
func WithTracer[T any](t Tracer) ActionOption[T] {
	return func(a *action[T]) {
		a.Tracer = t
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

// WithIOperation 设置IOperation
func WithIOperation[T any](o IOperation[T]) ActionOption[T] {
	return func(a *action[T]) {
		a.IOperation = o
	}
}

// WithIOperationX 设置IOperationX
func WithIOperationX[T any](o IOperationX[T]) ActionOption[T] {
	return func(a *action[T]) {
		a.IOperationX = o
	}
}

// WithIAssociation 设置IAssociation
func WithIAssociation[T any](o IAssociation) ActionOption[T] {
	return func(a *action[T]) {
		a.IAssociation = o
	}
}
