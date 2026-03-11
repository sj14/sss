package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestBasic(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("upload read-only", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "--read-only")
		must.Error(t, err)
		must.StrContains(t, err.Error(), "blocked by read-only mode")
	})

	t.Run("list after read-only upload attempt", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrNotContains(t, out, "README.md")
	})

	t.Run("upload 1", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md")
		must.NoError(t, err)
		must.StrContains(t, out, "README.md")
	})

	t.Run("upload 2", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "test/2")
		must.NoError(t, err)
		must.StrContains(t, out, "test/2")
	})

	t.Run("upload 3", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "test/3")
		must.NoError(t, err)
		must.StrContains(t, out, "test/3")
	})

	t.Run("list after upload", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  test/")
		must.StrContains(t, out, "README.md")
	})

	t.Run("list after upload with other delimiter", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls", "--delimiter=te")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  te")
		must.StrContains(t, out, "README.md")
	})

	t.Run("list after upload recursive", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls", "--delimiter=''")
		must.NoError(t, err)
		must.StrContains(t, out, "test/2")
		must.StrContains(t, out, "test/3")
		must.StrContains(t, out, "README.md")
	})

	t.Run("delete dry-run single", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "README.md", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "> dry-run mode <")
		must.StrContains(t, out, "README.md")
	})

	t.Run("delete dry-run dir no suffix", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "test", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "> dry-run mode <")
		must.StrContains(t, out, "failed to head object")
	})

	t.Run("delete dry-run dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "test/", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "> dry-run mode <")
		must.StrContains(t, out, "test/2")
		must.StrContains(t, out, "test/3")
	})

	t.Run("list after dry-run delete single", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  test/")
		must.StrContains(t, out, "README.md")
	})

	t.Run("list after dry-run delete dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls", "test/")
		must.NoError(t, err)
		must.StrContains(t, out, "2")
		must.StrContains(t, out, "3")
	})

	t.Run("delete single", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "README.md")
		must.NoError(t, err)
		must.StrContains(t, out, "README.md")
	})

	t.Run("list after delete single", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  test/")
		must.StrNotContains(t, out, "README.md")
	})

	t.Run("delete dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "test/")
		must.NoError(t, err)
		must.StrContains(t, out, "test/2")
		must.StrContains(t, out, "test/3")
	})

	t.Run("list after delete dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrNotContains(t, out, "test/")
		must.StrNotContains(t, out, "README.md")
	})
}
