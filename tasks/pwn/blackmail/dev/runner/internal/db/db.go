package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var ErrDuplicateApk = errors.New("duplicate apk")

type AppsDB struct {
	db *sql.DB
}

func NewAppsDB(db *sql.DB) *AppsDB {
	return &AppsDB{ db: db }
}

// Only call this from the 'submitter'.
func (r *AppsDB) EnsureSchema() error {
	_, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS apps(
	id	uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	ip	text NOT NULL,
	at	timestamp DEFAULT now(),
	hash text,
	apk bytea NOT NULL,
	pending boolean DEFAULT true,
	status text DEFAULT ''
);
CREATE INDEX IF NOT EXISTS apps_from_ip ON apps USING HASH(ip);

CREATE OR REPLACE FUNCTION apk_hash() RETURNS TRIGGER
	LANGUAGE plpgsql AS $$
BEGIN
	NEW.hash := substring(sha256(NEW.apk)::text FROM 3);
	RETURN NEW;
END;
$$;
CREATE OR REPLACE TRIGGER apk_hash_default BEFORE INSERT ON apps
	FOR EACH ROW EXECUTE FUNCTION apk_hash();
`)
	return err
}

func (r *AppsDB) Submit(ip string, apkContent []byte) (string, error) {
	q, err := r.db.Query(`INSERT INTO apps(ip, apk) VALUES ($1, $2) RETURNING id::text`, ip, apkContent)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return "", ErrDuplicateApk
		}
		return "", err
	}
	defer q.Close()
	
	var id string
	q.Scan(&id)
	return id, nil
}

func (r *AppsDB) GetLastAppTime(ip string) (time.Time, error) {
	q, err := r.db.Query(`SELECT at FROM apps WHERE ip=$1 ORDER BY at DESC LIMIT 1`, ip)
	if err != nil {
		return time.Time{}, err
	}
	defer q.Close()

	if q.Next() {
		var at time.Time
		q.Scan(&at)
		return at, nil
	}
	return time.Time{}, nil
}

type Submission struct {
	Id      string
	Status  string
	Pending bool
}

func (r *AppsDB) GetLastIPSubmission(ip string) (*Submission, error) {
	q, err := r.db.Query(`SELECT id, pending, status FROM apps WHERE ip=$1 ORDER BY at DESC LIMIT 1`, ip)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	if q.Next() {
		var s Submission
		q.Scan(&s.Id, &s.Pending, &s.Status)
		return &s, nil
	}
	return nil, nil
}

type PendingSubmission struct {
	Id	string
	Hash string
	Apk []byte
}

func (r *AppsDB) GetPendingSubmission() (*PendingSubmission, error) {
	q, err := r.db.Query(`SELECT id::text, hash, apk FROM apps WHERE pending = TRUE ORDER BY at ASC LIMIT 1`)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	if q.Next() {
		var ps PendingSubmission
		q.Scan(&ps.Id, &ps.Hash, &ps.Apk)
		return &ps, nil
	}
	return nil, nil
}

func (r *AppsDB) SetStatus(id string, pending bool, status string) error {
	q, err := r.db.Query(`UPDATE apps SET pending=$2, status=$3 WHERE id::text=$1`, id, pending, status)
	if err != nil {
		return err
	}
	return q.Close()
}
