package e2e

import (
	"context"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/sj14/sss/cli"
	"github.com/sj14/sss/util"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
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
	bucketName := util.RandomString(16, util.LettersLower)
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
