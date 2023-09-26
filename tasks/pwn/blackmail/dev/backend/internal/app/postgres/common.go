package postgres

type scanner interface {
	Scan(dest ...any) error
}
