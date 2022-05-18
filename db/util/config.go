package util

import "github.com/spf13/viper"

//Config stores all configuration of the application.
// The values are read by viper from a config file or enviroment variables

//here we will be using the unmarshalling feature of viper
type Config struct {
	//use variables declared in the app.env file
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	// tell viper the location of the location path
	viper.SetConfigName("app") // viper would look for file with this specific name
	viper.SetConfigType("env") // the type of the convig file
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
