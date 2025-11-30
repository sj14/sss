package main

import (
	"context"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/shoenig/test/must"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.IntN(len(letters))]
	}
	return string(s)
}

type safeWriter struct {
	sb strings.Builder
	mu sync.Mutex
}

func (w *safeWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.sb.Write(p)
}

var runMutex sync.Mutex

func run(ctx context.Context, args ...string) (string, error) {
	writer := &safeWriter{}

	runMutex.Lock()
	os.Args = append([]string{"sss"}, args...)
	err := exec(ctx, writer, writer)
	runMutex.Unlock()

	return writer.sb.String(), err
}

func createBucket(t *testing.T) string {
	bucketName := randomString(16)
	t.Attr("bucket", bucketName)

	t.Run("create bucket", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "mb")
		must.NoError(t, err, must.Sprint(out))
	})

	t.Cleanup(func() {
		// context from the test is already cancelled at that point
		ctx := context.Background()

		_, err := run(ctx, "bucket", bucketName, "rm", "/", "--force")
		must.NoError(t, err)

		_, err = run(ctx, "bucket", bucketName, "rb")
		must.NoError(t, err)
	})

	return bucketName
}

func TestBasic(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("upload read-only", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "put", "README.md", "--read-only")
		must.Error(t, err)
		must.StrContains(t, err.Error(), "blocked by read-only mode")
	})

	t.Run("list after read-only upload attempt", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrNotContains(t, out, "README.md")
	})

	t.Run("upload 1", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "README.md")
		must.NoError(t, err)
		must.StrContains(t, out, "README.md")
	})

	t.Run("upload 2", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "README.md", "test/2")
		must.NoError(t, err)
		must.StrContains(t, out, "test/2")
	})

	t.Run("upload 3", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "README.md", "test/3")
		must.NoError(t, err)
		must.StrContains(t, out, "test/3")
	})

	t.Run("list after upload", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "README.md")
		must.StrContains(t, out, "test/")
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
		must.StrContains(t, out, "test/")
		must.StrContains(t, out, "README.md")
	})

	t.Run("list after dry-run delete dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls", "test/")
		must.NoError(t, err)
		must.StrContains(t, out, "2")
		must.StrContains(t, out, "3")
	})

	t.Run("delete", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "README.md")
		must.NoError(t, err)
		must.StrContains(t, out, "README.md")
	})

	t.Run("list after delete", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "test/")
		must.StrNotContains(t, out, "README.md")
	})

	t.Run("delete dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "rm", "test/")
		must.NoError(t, err)
		must.StrContains(t, out, "2")
		must.StrContains(t, out, "3")
	})

	t.Run("list after delete dir", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls")
		must.NoError(t, err)
		must.StrNotContains(t, out, "test/")
		must.StrNotContains(t, out, "README.md")
	})
}

func TestSubdir(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("upload subdir 1 (dry-run)", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "util/", "test", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "test/util/progress/reader.go")
		must.StrContains(t, out, "test/util/zero.go")
	})

	t.Run("upload subdir 2", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "put", "util/", "test")
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
