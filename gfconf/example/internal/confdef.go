package internal

type (
	Application struct {
		App App `mapstructure:"app"`
		Log Log `mapstructure:"log"`
	}
	App struct {
		Mode string `mapstructure:"mode"`
		Port uint32 `mapstructure:"port"`
	}
	Log struct {
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"maxSize"`
		MaxBackups int    `mapstructure:"maxBackups"`
		MaxAge     int    `mapstructure:"maxAge"`
		Level      int    `mapstructure:"level"`
	}
)

type (
	WhiteAccount struct {
		Phone string `mapstructure:"phone"`
		ID    string `mapstructure:"id"`
	}
	Whitelist []WhiteAccount
)
