package config

type AppSettings struct {
	MongoDatabase struct {
		Uri      string `toml:"uri"`
		Database string `toml:"database"`
	} `toml:"mongo-database"`
	Redis struct {
		Address  string `toml:"address"`
		Password string `toml:"password"`
		DB       int    `toml:"db"`
		Protocol int    `toml:"protocol"`
	} `toml:"redis"`
	Auth struct {
		SecretKey       string `toml:"secret_key"`
		ExpirationHours int    `toml:"expiration_hours"`
		CookieKey       string `toml:"cookie_key"`
		Domain          string `toml:"domain"`
	} `toml:"auth"`
}
