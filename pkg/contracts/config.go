package contracts

type DBConfig interface {
	GetDBConfig() map[string]string
}

type AppConfig interface {
	GetAppConfig() map[string]string
}

type Config interface {
	AppConfig
	DBConfig
}
