package interfaces

import "context"

type MigrationService interface {
	Migrate(ctx context.Context) error
}
