package env

import "github.com/spf13/viper"

type Env struct {
	ImmichUrl          string
	ApiKey             string
	DbConnectionString string
	GotifyKey          string
	GotifyUrl          string
	OidcClientId       string
	OidcClientSecret   string
	OidcIssuerUrl      string
	OidcRedirectUrl    string
	ImmichApiKey       string
	AppPort            string
}

func New() *Env {
	env := &Env{}

	viper.BindEnv("IMMICH_URL")
	viper.SetDefault("IMMICH_URL", "https://immich.itsmelon.com")
	viper.BindEnv("IMMICH_API_KEY")

	viper.BindEnv("DB_CONNECTION_STRING")
	viper.SetDefault("DB_CONNECTION_STRING", "mongodb://root:password@localhost:27017")

	viper.BindEnv("GOTIFY_KEY")

	viper.BindEnv("GOTIFY_URL")
	viper.SetDefault("GOTIFY_URL", "https://gotify.itsmelon.com/message")

	viper.BindEnv("OIDC_CLIENT_ID")
	viper.BindEnv("OIDC_CLIENT_SECRET")
	viper.BindEnv("OIDC_ISSUER_URL")
	viper.BindEnv("OIDC_REDIRECT_URL")

	viper.BindEnv("API_KEY")

	viper.BindEnv("APP_PORT")
	viper.SetDefault("APP_PORT", "29442")

	env.ImmichUrl = viper.GetString("IMMICH_URl")
	env.ImmichApiKey = viper.GetString("IMMICH_API_KEY")
	env.ApiKey = viper.GetString("API_KEY")
	env.DbConnectionString = viper.GetString("DB_CONNECTION_STRING")
	env.GotifyKey = viper.GetString("GOTIFY_KEY")
	env.GotifyUrl = viper.GetString("GOTIFY_URL")
	env.OidcClientId = viper.GetString("OIDC_CLIENT_ID")
	env.OidcClientSecret = viper.GetString("OIDC_CLIENT_SECRET")
	env.OidcIssuerUrl = viper.GetString("OIDC_ISSUER_URL")
	env.OidcRedirectUrl = viper.GetString("OIDC_REDIRECT_URL")
	env.AppPort = viper.GetString("APP_PORT")
	return env
}
