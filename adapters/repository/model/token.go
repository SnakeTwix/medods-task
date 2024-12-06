package model

type Token struct {
	HashedRefreshTokenId string `gorm:"unique;not null"`
	TokenFamily          string `gorm:"primaryKey;not null"`
}
