package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestSubdir(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("upload subdir 1 (dry-run)", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../util/", "test", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "test/util/progress/reader.go")
		must.StrContains(t, out, "test/util/zero.go")
	})

	t.Run("upload subdir 2", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "../util/", "test")
		must.NoError(t, err)
		must.StrContains(t, out, "test/util/progress/reader.go")
		must.StrContains(t, out, "test/util/zero.go")
	})

	t.Run("download subdir 1 (dry-run)", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "test/", "yolo", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "yolo/test/util/progress/reader.go")
		must.StrContains(t, out, "yolo/test/util/zero.go")
	})

	t.Run("download subdir 2", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "get", "test/util/", "test/yolo", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "test/yolo/util/progress/reader.go")
		must.StrContains(t, out, "test/yolo/util/zero.go")
	})
}
