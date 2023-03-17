package domain

import (
	"errors"
	"fmt"
)

type Method struct {
	method string
}

var (
	MethodUnd        = Method{method: ""}            // "und"
	MethodDelete     = Method{method: "delete"}      // "delete"
	MethodMarkRead   = Method{method: "mark_read"}   // "mark_read"
	MethodMarkUnread = Method{method: "mark_unread"} // "mark_unread"
	MethodOrder      = Method{method: "order"}       // "order"
	MethodRemove     = Method{method: "remove"}      // "remove"
)

var ErrMethodSyntax = errors.New("unknown or unsupported method")

var stringsMethods = map[string]Method{
	MethodDelete.method:     MethodDelete,
	MethodMarkRead.method:   MethodMarkRead,
	MethodMarkUnread.method: MethodMarkUnread,
	MethodOrder.method:      MethodOrder,
	MethodRemove.method:     MethodRemove,
}

func ParseMethod(src string) (Method, error) {
	if method, ok := stringsMethods[src]; ok {
		return method, nil
	}

	return MethodUnd, fmt.Errorf("%w: %s", ErrMethodSyntax, src)
}

func (a Method) String() string {
	if a.method != "" {
		return a.method
	}

	return "und"
}
