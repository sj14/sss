package main_test

import (
	"math/rand"
	"os/exec"
	"testing"

	"github.com/shoenig/test/must"
)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func createBucket(t *testing.T) string {
	bucketName := randomString(16)
	t.Attr("bucket", bucketName)

	t.Run("create bucket", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "mb")
		_, err := cmd.CombinedOutput()
		must.NoError(t, err)
	})

	t.Cleanup(func() {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "rm", "/", "-force")
		err := cmd.Run()
		must.NoError(t, err)

		cmd = exec.Command("go", "run", "main.go", "-bucket", bucketName, "rb")
		err = cmd.Run()
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
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "put", "README.md", "-read-only")
		out, err := cmd.CombinedOutput()
		must.Error(t, err)
		must.StrContains(t, string(out), "blocked by read-only mode")
	})

	t.Run("list after read-only upload attempt", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "ls")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrNotContains(t, string(out), "README.md")
	})

	t.Run("upload 1", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "put", "README.md")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("upload 2", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "put", "README.md", "test/2")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "test/2")
	})

	t.Run("upload 3", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "put", "README.md", "test/3")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "test/3")
	})

	t.Run("list after upload", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "ls")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
		must.StrContains(t, string(out), "test/")
	})

	t.Run("delete dry-run single", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "rm", "README.md", "-dry-run")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "> dry-run mode <")
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("delete dry-run dir no suffix", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "rm", "test", "-dry-run")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "> dry-run mode <")
		must.StrContains(t, string(out), "failed to head object")
	})

	t.Run("delete dry-run dir", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "rm", "test/", "-dry-run")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "> dry-run mode <")
		must.StrContains(t, string(out), "test/2")
		must.StrContains(t, string(out), "test/3")
	})

	t.Run("list after dry-run delete single", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "ls")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "test/")
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("list after dry-run delete dir", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "ls", "test/")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "2")
		must.StrContains(t, string(out), "3")
	})

	t.Run("delete", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "rm", "README.md")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("list after delete", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "ls")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "test/")
		must.StrNotContains(t, string(out), "README.md")
	})

	t.Run("delete dir", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "rm", "test/")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "2")
		must.StrContains(t, string(out), "3")
	})

	t.Run("list after delete dir", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "ls")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrNotContains(t, string(out), "test/")
		must.StrNotContains(t, string(out), "README.md")
	})
}

func TestSubdir(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("upload subdir", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "put", "util", "test")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "test/util/progress/reader.go")
		must.StrContains(t, string(out), "test/util/zero.go")
	})

	t.Run("download subdir 1 (dry-run)", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "get", "test/", "yolo", "-dry-run")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "yolo/test/util/progress/reader.go")
		must.StrContains(t, string(out), "yolo/test/util/zero.go")
	})

	t.Run("download subdir 2", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "-bucket", bucketName, "get", "test/util/", "test/yolo", "-dry-run")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "test/yolo/util/progress/reader.go")
		must.StrContains(t, string(out), "test/yolo/util/zero.go")
	})
}
