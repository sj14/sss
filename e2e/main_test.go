package e2e

import (
	"context"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/sj14/sss/cli"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

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
	os.Args = append([]string{"sss", "--config=config.toml", "--profile=localstack"}, args...)
	err := cli.Exec(ctx, writer, writer, "e2e test")
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
