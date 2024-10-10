package settings

type Config struct {
	Server            ServerSetting         `mapstructure:"server"`
	PostgreSql        PostgreSqlSetting     `mapstructure:"postgreSql"`
	Logger            LoggerSetting         `mapstructure:"logger"`
	Redis             RedisSetting          `mapstructure:"redis"`
	Authentication    AuthenticationSetting `mapstructure:"authentication"`
	CloudinarySetting CloudinarySetting     `mapstructure:"cloudinary"`
}

type ServerSetting struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type PostgreSqlSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	SslMode         string `mapstructure:"sslMode"`
}

type RedisSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type LoggerSetting struct {
	LogLevel    string `mapstructure:"log_level"`
	FileLogName string `mapstructure:"file_log_name"`
	MaxBackups  int    `mapstructure:"max_backup"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxSize     int    `mapstructure:"max_size"`
	Compress    bool   `mapstructure:"compress"`
}

type AuthenticationSetting struct {
	JwtScretKey string `mapstructure:"jwtScretKey"`
}

type CloudinarySetting struct {
	CloudName    string `mapstructure:"cloud_name"`
	ApiKey       string `mapstructure:"api_key"`
	ApiSecretKey string `mapstructure:"api_secret_key"`
	Folder       string `mapstructure:"folder"`
}
