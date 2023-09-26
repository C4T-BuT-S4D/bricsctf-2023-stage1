// Domain definitions
package app

type UserService interface {
	Register(username, password string) (*User, error)
	GetById(id Uid) (*User, error)
	GetByName(name string) (*User, error)
	GetBots() ([]User, error)
	GetByCreds(username, password string) (*User, error)
	UserExists(id Uid) bool
}

type Uid int

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Id           Uid    `json:"id"`
	IsBot        bool   `json:"-"`
}

type ChatsService interface {
	GetChat(me, other Uid) (*Chat, error)
	// Returns slice of chats with latest message only
	GetChatsPreview(me Uid) ([]Chat, error)
	SendMessage(me, other Uid, msg Message) error
}

const (
	MessageTypeText    = "text"
	MessageTypeSticker = "sticker"
)

type Chat struct {
	Messages []Message `json:"messages"`
	With     Uid       `json:"with"`
}

type Message struct {
	Text    *string  `json:"text"`
	Sticker *Sticker `json:"sticker"`
	Type    string   `json:"type"`
	From    Uid      `json:"from"`
}

type Sticker struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
