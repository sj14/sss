package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestPut(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("upload single file", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "README.md")
	})

	t.Run("upload single file in dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "mydir/", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "mydir/README.md")
	})

	t.Run("upload single file with specific name", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "my-file", "--dry-run")
		must.NoError(t, err)
		must.StrNotContains(t, out, "README.md")
		must.StrContains(t, out, "my-file")
	})

	t.Run("upload single file with specific name in dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "mydir/my-file", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "mydir/my-file")
	})

	t.Run("upload dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../util", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "util/zero.go")
		must.StrContains(t, out, "util/progress/reader.go")
	})

	t.Run("upload dir into subdir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../util", "mydir", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "mydir/util/zero.go")
		must.StrContains(t, out, "mydir/util/progress/reader.go")
	})
}
