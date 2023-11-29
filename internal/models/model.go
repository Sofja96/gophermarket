package models

import "time"

const (
	ContextKeyUser = "login"
)

type OrderStatus = string

const (
	NEW        OrderStatus = "NEW"
	INVALID    OrderStatus = "INVALID"    //заказ не принят к расчёту, и вознаграждение не будет начислено
	PROCESSING OrderStatus = "PROCESSING" //расчёт начисления в процессе
	PROCESSED  OrderStatus = "PROCESSED"  //расчёт начисления окончен
	REGISTERED OrderStatus = "REGISTERED" //заказ зарегистрирован, но вознаграждение не рассчитано
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type UserWithdrawal struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}
