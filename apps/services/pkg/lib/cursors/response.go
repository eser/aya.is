package cursors

type Cursored[T any] struct {
	Data      T       `json:"data"`
	CursorPtr *string `json:"cursor"`
}

func WrapResponseWithCursor[T any](data T, cursorPtr *string) Cursored[T] {
	return Cursored[T]{
		Data:      data,
		CursorPtr: cursorPtr,
	}
}
