package port

type ConfigService interface {
	GetAccessTokenSignKey() []byte
	GetRefreshTokenSignKey() []byte
	GetRefreshTokenHashKey() []byte
}
