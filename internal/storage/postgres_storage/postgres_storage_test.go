package postgres_storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gobserver/internal/entity"
	"testing"
	"time"
)

func TestPostgresStorage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		ps, err := NewPostgresStorage()
		assert.NoError(t, err)

		_, err = ps.AddFileChange(context.Background(), entity.FileChange{
			ProjectPath: "biba/boba",
			FilePath:    "biba/boba/AHAHAAHAH_PANIC.txt",
			Operation:   "CREATE",
			Timestamp:   time.Now(),
		})
		assert.NoError(t, err)

		_, err = ps.AddCommandExecution(context.Background(), entity.CommandExecution{
			ProjectPath: "bebra",
			Command:     "rm -rf /*",
			ExitCode:    1,
			Timestamp:   time.Now().Add(10 * 365 * 24 * time.Hour),
		})
		assert.NoError(t, err)
	})
}
