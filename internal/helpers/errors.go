package helpers

import "errors"

// var ErrWrongOrderNumber = errors.New("wrong order number")
var ErrAnotherUserOrder = errors.New("order number already exists for another user")
var ErrExistsOrder = errors.New("order number already exists")

var ErrInsufficientBalance = errors.New("insufficient balance")

//var ErrNotEnoughBalance = errors.New("not enough balance")
