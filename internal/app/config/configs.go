package config

type Config struct {
	Address        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseDSN    string `env:"DATABASE_URI"`
}

func LoadConfig() *Config {
	var cfg Config
	return &cfg

}
