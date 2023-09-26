package postgres

import (
	"database/sql"
	"errors"

	"cbs.dev/brics/droidchat/internal/app"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Implements app.UserService
type UserService struct {
	db *sql.DB
}

// Create new postgres app.UserService and init tables.
// Must precede NewChatService
func NewUserService(db *sql.DB) app.UserService {
	s := UserService{db: db}
	s.ensureSchema()
	return &s
}

func (u *UserService) ensureSchema() {
	_, err := u.db.Exec(`
CREATE TABLE IF NOT EXISTS "user" (
	id serial PRIMARY KEY,
	username text UNIQUE NOT NULL,
	password_hash text NOT NULL,
	is_bot	bool DEFAULT false
);
`)
	if err != nil {
		// catastrophe
		panic(errors.Join(errors.New("user schema failed"), err))
	}
}

// GetByCreds implements app.UserService
func (s *UserService) GetByCreds(username string, password string) (*app.User, error) {
	user, err := s.GetByName(username)
	if err != nil {
		return nil, errors.Join(err, app.ErrWrongCreds)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, app.ErrWrongCreds
	}
	return user, nil
}

// GetBots implements app.UserService
func (s *UserService) GetBots() ([]app.User, error) {
	return s.selectBotUsers()
}

func (s *UserService) selectBotUsers() ([]app.User, error) {
	rows, err := s.db.Query(`SELECT * FROM "user" WHERE is_bot`)
	defer rows.Close()
	if err != nil {
		panic(err) // idk
	}
	users := make([]app.User, 0)
	for rows.Next() {
		user := mustScanUser(rows)
		users = append(users, user)
	}
	return users, nil
}

// GetById implements app.UserService
func (s *UserService) GetById(id app.Uid) (*app.User, error) {
	if user := s.selectUserById(id); user == nil {
		return nil, app.ErrNotFound
	} else {
		return user, nil
	}
}

func (s *UserService) selectUserById(id app.Uid) *app.User {
	rows, err := s.db.Query(`SELECT * FROM "user" WHERE id = $1`, id)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	for rows.Next() {
		u := mustScanUser(rows)
		return &u
	}

	return nil
}

// GetByName implements app.UserService
func (s *UserService) GetByName(name string) (*app.User, error) {
	if user := s.selectUserByName(name); user == nil {
		return nil, app.ErrNotFound
	} else {
		return user, nil
	}
}

func (s *UserService) selectUserByName(name string) *app.User {
	rows, err := s.db.Query(`SELECT * FROM "user" WHERE username = $1`, name)
	defer rows.Close()
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}

	for rows.Next() {
		u := mustScanUser(rows)
		return &u
	}

	return nil
}

// Register implements app.UserService
func (s *UserService) Register(username string, password string) (*app.User, error) {
	newUser := &app.User{
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if err := s.insertUser(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func hashPassword(password string) string {
	// ignoring err because gin handler validates input
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h)
}

// inserts user into db and updates user.Id
func (u *UserService) insertUser(user *app.User) error {
	row := u.db.QueryRow(
		`INSERT INTO "user"(username, password_hash) VALUES ($1, $2) RETURNING id`,
		user.Username,
		user.PasswordHash,
	)

	var newId int
	err := row.Scan(&newId)
	if err == nil {
		user.Id = app.Uid(newId)
		return nil
	}

	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return app.ErrExists
	} else {
		panic(err) // idk
	}
}

// UserExists implements app.UserService
func (s *UserService) UserExists(id app.Uid) bool {
	row := s.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM "user" WHERE id = $1)`, id)
	var result bool
	if err := row.Scan(&result); err != nil {
		panic(err)
	}
	return result
}

// Read *User from sql.Row/sql.Rows (helper)
func mustScanUser(sc *sql.Rows) app.User {
	var user app.User
	if err := sc.Scan(
		&user.Id,
		&user.Username,
		&user.PasswordHash,
		&user.IsBot,
	); err != nil {
		panic(err) // catastrophe
	}
	return user
}
