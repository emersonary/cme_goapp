package config

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

var cfg *Conf

type Conf struct {
	WaitFroStartUp int    `mapstructure:"APP_WAITFORSTARTUP"`
	DBDriver       string `mapstructure:"DB_DRIVER"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         int    `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_DBUSER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	WebServerPort  int    `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret      string `mapstructure:"WEB_JWT_SECRET"`
	JWTExpiresIn   int    `mapstructure:"WEB_JWT_EXPIRESIN"`
	RedisHost      string `mapstructure:"REDIS_HOST"`
	RedisPort      int    `mapstructure:"REDIS_PORT"`
	TokenAuth      *jwtauth.JWTAuth
}

func LoadConfig(path string) (*Conf, error) {
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, err
}

// funcao init eh executada antes do main
/*
func init() {

}
*/
