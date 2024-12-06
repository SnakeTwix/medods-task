package configsrv

import "medods-api/util/env"

type ConfigService struct {
	accessTokenSignKey  []byte
	refreshTokenSignKey []byte
	refreshTokenHashKey []byte
}

func New() *ConfigService {
	// In this case, we should initialize the config right away from the env variables
	// Don't want the application randomly crashing because an env variable was not defined
	return &ConfigService{
		accessTokenSignKey:  []byte(env.Get(env.ACCESS_TOKEN_SIGN_KEY)),
		refreshTokenSignKey: []byte(env.Get(env.REFRESH_TOKEN_SIGN_KEY)),
		refreshTokenHashKey: []byte(env.Get(env.REFRESH_TOKEN_HASH_KEY)),
	}
}

func (c *ConfigService) GetAccessTokenSignKey() []byte {
	return c.accessTokenSignKey
}
func (c *ConfigService) GetRefreshTokenSignKey() []byte {
	return c.refreshTokenSignKey

}
func (c *ConfigService) GetRefreshTokenHashKey() []byte {
	return c.refreshTokenHashKey

}
