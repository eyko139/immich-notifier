package env

import "github.com/spf13/viper"

type Env struct {
	ImmichUrl          string
	ApiKey             string
	DbConnectionString string
	GotifyKey string
}

func New() *Env {
	env := &Env{}

	viper.BindEnv("IMMICH_URL")
	viper.SetDefault("IMMICH_URL", "https://immich.itsmelon.com")
	viper.BindEnv("DB_CONNECTION_STRING")
	viper.SetDefault("DB_CONNECTION_STRING", "mongodb://root:password@localhost:27017")
	viper.BindEnv("GOTIFY_KEY")

	viper.BindEnv("API_KEY")

	env.ImmichUrl = viper.GetString("IMMICH_URl")
	env.ApiKey = viper.GetString("API_KEY")
	env.DbConnectionString = viper.GetString("DB_CONNECTION_STRING")
	env.GotifyKey = viper.GetString("GOTIFY_KEY")
	return env
}
