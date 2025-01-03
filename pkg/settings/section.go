package settings

type Config struct {
	Server         ServerSetting         `mapstructure:"server"`
	PostgreSql     PostgreSqlSetting     `mapstructure:"postgresql"`
	Logger         LoggerSetting         `mapstructure:"logger"`
	Media          MediaSetting          `mapstructure:"media"`
	Redis          RedisSetting          `mapstructure:"redis"`
	Authentication AuthenticationSetting `mapstructure:"authentication"`
	MailService    MailServiceSetting    `mapstructure:"mail_service"`
	MomoSetting    MomoSetting           `mapstructure:"momo"`
	GoogleSetting  GoogleSetting         `mapstructure:"google"`
}

type ServerSetting struct {
	Port           int    `mapstructure:"port"`
	Mode           string `mapstructure:"mode"`
	ServerEndpoint string `mapstructure:"server_endpoint"`
}

type PostgreSqlSetting struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Dbname          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	SslMode         string `mapstructure:"ssl_mode"`
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
	JwtSecretKey      string `mapstructure:"jwtSecretKey"`
	JwtAdminSecretKey string `mapstructure:"jwtAdminSecretKey"`
}

type MediaSetting struct {
	Folder string `mapstructure:"folder"`
}

type MailServiceSetting struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     string `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
}

type MomoSetting struct {
	PartnerCode  string `mapstructure:"partner_code"`
	AccessKey    string `mapstructure:"access_key"`
	SecretKey    string `mapstructure:"secret_key"`
	RedirectUrl  string `mapstructure:"redirect_url"`
	IpnURL       string `mapstructure:"ipn_url"`
	EndpointHost string `mapstructure:"endpoint_host"`
	EndpointPath string `mapstructure:"endpoint_path"`
}

type GoogleSetting struct {
	GoogleTokensUrl string `mapstructure:"google_tokens_url"`
	SecretId        string `mapstructure:"secret_id"`
	WebClientId     string `mapstructure:"web_client_id"`
	AndroidClientId string `mapstructure:"android_client_id"`
	IosClientId     string `mapstructure:"ios_client_id"`
	RedirectUrl     string `mapstructure:"redirect_url"`
}
