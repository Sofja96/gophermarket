package helpers

import "errors"

var ErrAnotherUserOrder = errors.New("order number already exists for another user")
var ErrExistsOrder = errors.New("order number already exists")

var ErrInsufficientBalance = errors.New("insufficient balance")
