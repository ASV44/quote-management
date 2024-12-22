package config

type Config struct {
	DbURL string `env:"DB_URL,required"`
	Addr  string `env:"API_HTTP_PORT,required"`
}
