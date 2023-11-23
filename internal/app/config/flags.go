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
	flag.StringVar(&s.Address, "a", "", "address and port to run server")
	flag.StringVar(&s.AccrualAddress, "r", "", "accrual system address ")
	flag.StringVar(&s.DatabaseDSN, "d", "", "connect to database")
	//flag.StringVar(&s.HashKey, "k", "", "key for hash")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		s.Address = envRunAddr
	}
	if envRunAddr := os.Getenv("DATABASE_DSN"); envRunAddr != "" {
		s.DatabaseDSN = envRunAddr
	}
	if envRunAddr := os.Getenv("DATABASE_DSN"); envRunAddr != "" {
		s.AccrualAddress = envRunAddr
	}
	//if envRunAddr := os.Getenv("KEY"); envRunAddr != "" {
	//	s.HashKey = envRunAddr
	//}
}
