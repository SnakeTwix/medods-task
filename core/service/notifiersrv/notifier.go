package notifiersrv

import (
	"context"
	"github.com/google/uuid"
	"log"
)

type NotifierService struct {
}

func New() *NotifierService {
	return &NotifierService{}
}

func (s *NotifierService) NotifyUserIPChange(context context.Context, userId uuid.UUID, newIp string) error {
	// Would have sent an email here

	log.Printf("OH NO, user with id: %s changed their ip to: %s\n", userId, newIp)

	return nil
}
