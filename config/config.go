package config

type AppSettings struct {
	MongoDatabase struct {
		Uri string `toml:"uri"`
	} `toml:"mongo-database"`
}
