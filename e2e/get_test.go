package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestGet(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("prepare", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../util")
		must.NoError(t, err)
		must.StrContains(t, out, "util/zero.go")
		must.StrContains(t, out, "util/progress/reader.go")
	})

	t.Run("get single file", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "util/zero.go", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "util/zero.go")
	})

	t.Run("get single file to specific dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "util/zero.go", "my-dir/", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "my-dir/zero.go")
	})

	t.Run("get single file to specific name", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "util/zero.go", "my-file", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "my-file")
	})

	t.Run("get single file to specific dir with specific name", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "util/zero.go", "my-dir/my-file", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "my-dir/my-file")
	})

	t.Run("get dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "util/", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "util/zero.go")
		must.StrContains(t, out, "util/progress/reader.go")
	})

	t.Run("get dir to specific dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "util/", "mydir", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "mydir/util/zero.go")
		must.StrContains(t, out, "mydir/util/progress/reader.go")
	})
}
