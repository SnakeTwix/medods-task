package configsrv

import "medods-api/util/env"

type ConfigService struct {
	accessTokenSignKey  []byte
	refreshTokenSignKey []byte
	refreshTokenHashKey []byte
}

func New() *ConfigService {
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
