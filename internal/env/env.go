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
	AppEnv             string
	ImmichPollInterval int
	WebsiteURL         string
	BotURL             string
	AppVersion         string
}

func New() *Env {
	env := &Env{}

	viper.BindEnv("IMMICH_URL")
	viper.SetDefault("IMMICH_URL", "https://immich.itsmelon.com")
	viper.BindEnv("IMMICH_API_KEY")

	err := viper.BindEnv("DB_CONNECTION_STRING")
	if err != nil {
		panic(err)
	}

	viper.BindEnv("GOTIFY_KEY")

	viper.BindEnv("GOTIFY_URL")
	viper.SetDefault("GOTIFY_URL", "https://gotify.itsmelon.com/message")

	viper.BindEnv("WEBSITE_URL")
	viper.SetDefault("WEBSITE_URL", "https://bot.itsmelon.com")

	viper.BindEnv("BOT_URL")
	viper.SetDefault("BOT_URL", "https://api.telegram.org/bot6429398075:AAFjoY4mthOBReLML8qh_-Zj_K9LZdKWQKc")

	viper.BindEnv("IMMICH_POLL_INTERVAL_SECONDS")
	viper.SetDefault("IMMICH_POLL_INTERVAL_SECONDS", 60)

	viper.BindEnv("OIDC_CLIENT_ID")
	viper.BindEnv("OIDC_CLIENT_SECRET")
	viper.BindEnv("OIDC_ISSUER_URL")
	viper.BindEnv("OIDC_REDIRECT_URL")

	viper.BindEnv("API_KEY")

	viper.BindEnv("APP_VERSION")

	viper.BindEnv("APP_PORT")
	viper.SetDefault("APP_PORT", "29442")

	viper.BindEnv("APP_ENV")
	viper.SetDefault("APP_ENV", "development")

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
	env.AppEnv = viper.GetString("APP_ENV")
	env.ImmichPollInterval = viper.GetInt("IMMICH_POLL_INTERVAL_SECONDS")
	env.BotURL = viper.GetString("BOT_URL")
    env.AppVersion = viper.GetString("APP_VERSION")
	return env
}
