package common

// CommonConfig описывает общие настройки для всех сервисов.
type CommonConfig struct {

	// Database содержит настройки подключения к базе данных.
	Database struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"database"`

	// Параметры запуска HTTP сервер
	Http struct {
		ReadTimeout  int `mapstructure:"read_timeout"`
		WriteTimeout int `mapstructure:"write_timeout"`
	} `mapstructure:"http"`
}

// ServiceConfig содержит настройки специфичные для сервисов.
type ServiceConfig struct {
	CommonConfig `mapstructure:",squash"`

	// Address содержит настройки сетевого адреса сервиса.
	Address struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"address"`
}
