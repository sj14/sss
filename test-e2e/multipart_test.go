package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestMultipart(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("create in root", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "multipart", "create", "yolo")
		must.NoError(t, err)
	})

	t.Run("create in subdir/prefix", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "multipart", "create", "mydir/something")
		must.NoError(t, err)
	})

	t.Run("list", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  mydir/")
		must.StrContains(t, out, "yolo")
	})

	t.Run("list dir/prefix", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls", "mydir/")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list resursive", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls", "-d ''")
		must.NoError(t, err)
		must.StrContains(t, out, "mydir/something")
		must.StrContains(t, out, "yolo")
	})
}
