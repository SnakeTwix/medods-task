package port

import (
	"context"
	"github.com/google/uuid"
)

type NotifierService interface {
	NotifyUserIPChange(context context.Context, userId uuid.UUID, newIp string) error
}
