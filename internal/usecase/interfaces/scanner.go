package interfaces

type Scanner interface {
	Scan(dest ...any) error
}
