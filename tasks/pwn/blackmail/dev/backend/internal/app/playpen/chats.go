// This package implements app.ChatsService and restrictions, namely:
// 1. Users can only chat with bot users.
// 2. Flag user is not allowed to interact with anyone.
package playpen

import (
	"fmt"

	"cbs.dev/brics/droidchat/internal/app"
)

var ErrPlaypenForbidden error = fmt.Errorf("forbidden")

type SafeChatsService struct {
	us app.UserService
	// intentionally public to allow bypassing
	Cs app.ChatsService
}

func NewChatsService(us app.UserService, cs app.ChatsService) app.ChatsService {
	return &SafeChatsService{us: us, Cs: cs}
}

/*
GetChat iements app.ChatsService
*/
func (s *SafeChatsService) GetChat(me app.Uid, other app.Uid) (*app.Chat, error) {
	chat, err := s.Cs.GetChat(me, other)
	if err != nil {
		return nil, err
	}

	{
		// Assumed
		userMe, _ := s.us.GetById(me)
		// Error checked earlier
		userOther, _ := s.us.GetById(other)
		// Both real users
		if !userMe.IsBot && !userOther.IsBot {
			return nil, ErrPlaypenForbidden
		}
	}

	return chat, nil
}

/*
GetChatsPreview implements app.ChatsService
*/
func (s *SafeChatsService) GetChatsPreview(me app.Uid) ([]app.Chat, error) {
	// No restrictions
	return s.Cs.GetChatsPreview(me)
}

/*
SendMessage implements app.ChatsService
*/
func (s *SafeChatsService) SendMessage(me app.Uid, other app.Uid, msg app.Message) error {
	{
		userMe, err := s.us.GetById(me)
		if err != nil {
			return err
		}
		if userMe.Username == "admin" {
			return ErrPlaypenForbidden
		}
		userOther, err := s.us.GetById(other)
		if err != nil {
			return err
		}
		// Both real users
		if !userMe.IsBot && !userOther.IsBot {
			return ErrPlaypenForbidden
		}
	}

	return s.Cs.SendMessage(me, other, msg)
}
