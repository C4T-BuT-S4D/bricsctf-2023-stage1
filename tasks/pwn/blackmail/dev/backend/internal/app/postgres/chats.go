package postgres

import (
	"database/sql"
	"fmt"

	"cbs.dev/brics/droidchat/internal/app"
	"cbs.dev/brics/droidchat/internal/app/config"
)

// Implements app.ChatsService
type ChatsService struct {
	db *sql.DB
	u  app.UserService
}

func NewChatsService(db *sql.DB, u app.UserService) app.ChatsService {
	s := ChatsService{db: db, u: u}
	s.ensureSchema()
	return &s
}

func (s *ChatsService) ensureSchema() {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS message (
	uid_from int REFERENCES "user",
	uid_to int REFERENCES "user",
	at timestamp DEFAULT now(),
	msg varchar(1024),
	sticker_id varchar(64)
);
CREATE INDEX IF NOT EXISTS message_from_to_idx ON message(uid_from, uid_to);
CREATE INDEX IF NOT EXISTS message_from_idx ON message(uid_from);
CREATE INDEX IF NOT EXISTS message_to_idx ON message(uid_to);
`)
	if err != nil {
		panic(fmt.Sprintf("message schema failed: %v", err))
	}
}

// GetChat implements app.ChatsService
func (s *ChatsService) GetChat(me app.Uid, other app.Uid) (*app.Chat, error) {
	if _, err := s.u.GetById(other); err != nil {
		return nil, err
	}
	// We assume that 'me' is always valid

	return s.getChat(me, other), nil
}

func (s *ChatsService) getChat(me, other app.Uid) *app.Chat {
	rows, err := s.db.Query(`
SELECT uid_from, COALESCE(msg, ''), COALESCE(sticker_id, '') FROM message 
	WHERE uid_from = $1 AND uid_to = $2 
	OR uid_from = $2 AND uid_to = $1
	ORDER BY at
`, me, other)
	defer rows.Close()
	if err != nil {
		panic(err) // not related to user input
	}

	messages := make([]app.Message, 0)
	for rows.Next() {
		messages = append(messages, mustScanMessage(rows))
	}

	return &app.Chat{
		With:     other,
		Messages: messages,
	}
}

// GetChatsPreview implements app.ChatsService
func (s *ChatsService) GetChatsPreview(me app.Uid) ([]app.Chat, error) {
	return s.GetChats(me), nil
}

func (s *ChatsService) GetChats(user app.Uid) []app.Chat {
	// TODO: Optimize the query?
	rows, err := s.db.Query(`
WITH mychats AS (
	SELECT uid_from, uid_to FROM message
	WHERE uid_from = $1 OR uid_to = $1
),
chats AS (
	SELECT DISTINCT ON (CASE WHEN uid_from = $1 THEN uid_to ELSE uid_from END)
		(CASE WHEN uid_from = $1 THEN uid_to ELSE uid_from END) AS other_id
	FROM mychats
)
SELECT (CASE WHEN latest.uid_from = $1 THEN latest.uid_to ELSE latest.uid_from END), latest.uid_from, latest.msg, latest.sticker_id FROM chats
JOIN LATERAL (
	SELECT uid_from, uid_to, COALESCE(msg, '') as msg, COALESCE(sticker_id, '') as sticker_id FROM message
	WHERE uid_from = chats.other_id AND uid_to = $1 
		OR uid_from = $1 AND uid_to = chats.other_id
	ORDER BY at DESC
	LIMIT 1
) latest ON latest.uid_from = chats.other_id OR latest.uid_to = chats.other_id
	`, user)
	defer rows.Close()
	if err != nil {
		panic(err) // not related to user input
	}

	chats := make([]app.Chat, 0)
	for rows.Next() {
		chats = append(chats, mustScanPreview(rows))
	}
	return chats
}

// SendMessage implements app.ChatsService
func (s *ChatsService) SendMessage(me app.Uid, other app.Uid, msg app.Message) error {
	if _, err := s.u.GetById(other); err != nil {
		return err
	}
	// We assume that 'me' is always valid

	msg.From = me
	s.insertMessage(msg, other)
	return nil
}

func (s *ChatsService) insertMessage(msg app.Message, to app.Uid) {
	var err error
	var r *sql.Rows
	if msg.Type == app.MessageTypeText {
		r, err = s.db.Query(`
INSERT INTO message(uid_from, uid_to, msg)
VALUES ($1, $2, $3)
		`, msg.From, to, msg.Text)
	} else {
		r, err = s.db.Query(`
INSERT INTO message(uid_from, uid_to, sticker_id)
VALUES ($1, $2, $3)
		`, msg.From, to, msg.Sticker.Id)
	}
	defer r.Close()
	if err != nil {
		panic(err) // not related to user input
	}
}

// Helper function
func mustScanMessage(sc *sql.Rows) app.Message {
	var msg app.Message

	var text string
	var stickerId string
	if err := sc.Scan(
		&msg.From,
		&text,
		&stickerId,
	); err != nil {
		panic(err) // prolly a parse error
	}

	if text != "" {
		msg.Type = app.MessageTypeText
		msg.Text = &text
	} else {
		msg.Type = app.MessageTypeSticker
		msg.Sticker = config.AvailableStickers[stickerId]
	}

	return msg
}

func mustScanPreview(sc *sql.Rows) app.Chat {
	var with app.Uid
	var text string
	var stickerId string
	var msg app.Message

	if err := sc.Scan(
		&with,
		&msg.From,
		&text,
		&stickerId,
	); err != nil {
		panic(err) // no idea
	}

	if text != "" {
		msg.Type = app.MessageTypeText
		msg.Text = &text
	} else {
		msg.Type = app.MessageTypeSticker
		msg.Sticker = config.AvailableStickers[stickerId]
	}

	return app.Chat{
		With:     with,
		Messages: []app.Message{msg},
	}
}
