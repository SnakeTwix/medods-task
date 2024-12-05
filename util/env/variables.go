package env

// Type internal type signifies an Environment variable
type Type string

const (
	POSTGRES_DB            Type = "POSTGRES_DB"
	POSTGRES_USER          Type = "POSTGRES_USER"
	POSTGRES_PASSWORD      Type = "POSTGRES_PASSWORD"
	API_PORT               Type = "API_PORT"
	ACCESS_TOKEN_SIGN_KEY  Type = "ACCESS_TOKEN_SIGN_KEY"
	REFRESH_TOKEN_SIGN_KEY Type = "REFRESH_TOKEN_SIGN_KEY"
	REFRESH_TOKEN_HASH_KEY Type = "REFRESH_TOKEN_HASH_KEY"
)
