package repo

type Scanner interface {
	Scan(dest ...interface{}) error
}
