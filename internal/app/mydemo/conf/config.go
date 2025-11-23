package conf

// Config 配置yaml结构
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"database"`

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		CacheTTL int    `yaml:"cache_ttl"`
	} `yaml:"redis"`

	Class []string `yaml:"class"`

	WeatherAppKey   string `yaml:"weather_app_key"`
	WeatherSignSalt string `yaml:"weather_sign_salt"`
}
