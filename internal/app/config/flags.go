package config

import "C"
import (
	"flag"
	"os"
)

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags(s *Config) {
	//var conf Config
	flag.StringVar(&s.Address, "a", "localhost:8081", "address and port to run server")
	flag.StringVar(&s.AccrualAddress, "r", "http://localhost:8080", "accrual system address")
	flag.StringVar(&s.DatabaseDSN, "d", "postgres://gophermarket:userpassword@localhost:5432/gophermarket?sslmode=disable", "connect to database")
	//flag.StringVar(&s.HashKey, "k", "", "key for hash")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		s.Address = envRunAddr
	}
	if envRunAddr := os.Getenv("DATABASE_DSN"); envRunAddr != "" {
		s.DatabaseDSN = envRunAddr
	}
	if envRunAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envRunAddr != "" {
		s.AccrualAddress = envRunAddr
	}
	//if envRunAddr := os.Getenv("KEY"); envRunAddr != "" {
	//	s.HashKey = envRunAddr
	//}
}
