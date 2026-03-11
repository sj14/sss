package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestCleanupMultipart(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("create mp in root", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "multipart", "create", "something")
		must.NoError(t, err)
	})

	t.Run("create mp in subdir/prefix", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "multipart", "create", "directory/something")
		must.NoError(t, err)
	})

	t.Run("list mp", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  directory/")
		must.StrContains(t, out, "something")
	})

	t.Run("cleanup dry-run", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "cleanup", "--force", "--all-multiparts", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list mp after dry-run", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  directory/")
		must.StrContains(t, out, "something")
	})

	t.Run("cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "cleanup", "--force", "--all-multiparts")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list mp after cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrNotContains(t, out, "PREFIX  directory/")
		must.StrNotContains(t, out, "something")
	})
}
