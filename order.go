package query

var _ IOrder[any] = (*Order[any])(nil)

type (
	IOrder[T any] interface {
		Desc() IAction[T]
		Asc() IAction[T]
	}

	Order[T any] struct {
		column string
		asc    orderType
		IAction[T]
	}

	orderType bool
)

const (
	ASC  orderType = true
	DESC orderType = false
)

func (o orderType) String() string {
	if o {
		return "ASC"
	}
	return "DESC"
}

// NewOrder 实例化排序
func NewOrder[T any](column string) *Order[T] {
	return &Order[T]{
		column: column,
	}
}

func (o *Order[T]) Desc() IAction[T] {
	return o.IAction.WithDB(o.IAction.DB().Order("`" + o.column + "` " + ASC.String()))
}

func (o *Order[T]) Asc() IAction[T] {
	return o.IAction.WithDB(o.IAction.DB().Order("`" + o.column + "` " + DESC.String()))
}

func (o *Order[T]) WithIAction(action IAction[T]) *Order[T] {
	o.IAction = action
	return o
}
