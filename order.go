package query

var _ IOrder[any] = (*Order[any])(nil)

type (
	IOrder[T any] interface {
		Desc() IAction[T]
		Asc() IAction[T]
		Column() string
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

func NewOrder[T any](column string) *Order[T] {
	return &Order[T]{
		column: column,
	}
}

func (o *Order[T]) Desc() IAction[T] {
	o.asc = false
	return o.IAction.WithDB(o.IAction.DB().Order("`" + o.column + "` " + o.asc.String()))
}

func (o *Order[T]) Asc() IAction[T] {
	o.asc = true

	return o.IAction.WithDB(o.IAction.DB().Order("`" + o.column + "` " + o.asc.String()))
}

func (o *Order[T]) Column() string {
	return o.column
}

func (o *Order[T]) WithIAction(action IAction[T]) *Order[T] {
	o.IAction = action
	return o
}
