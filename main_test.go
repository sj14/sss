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

func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := randomString(16)
	t.Attr("bucket", bucketName)

	t.Cleanup(func() {
		cmd := exec.Command("go", "run", "main.go", "rb", "-bucket", bucketName)
		err := cmd.Run()
		must.NoError(t, err)
	})

	t.Run("create bucket", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "mb", "-bucket", bucketName)
		_, err := cmd.CombinedOutput()
		must.NoError(t, err)
	})

	t.Run("upload read-only", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "put", "-bucket", bucketName, "README.md", "-read-only")
		out, err := cmd.CombinedOutput()
		must.Error(t, err)
		must.StrContains(t, string(out), "read-only mode: blocked PUT")
	})

	t.Run("list after read-only upload attempt", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "ls", "-bucket", bucketName)
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrNotContains(t, string(out), "README.md")
	})

	t.Run("upload", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "put", "-bucket", bucketName, "README.md")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("list after upload", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "ls", "-bucket", bucketName)
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("delete dry-run", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "rm", "-bucket", bucketName, "README.md", "-dry-run")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "> dry-run mode <")
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("list after dry-run delete", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "ls", "-bucket", bucketName)
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("delete", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "rm", "-bucket", bucketName, "README.md")
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrContains(t, string(out), "README.md")
	})

	t.Run("list after delete", func(t *testing.T) {
		cmd := exec.Command("go", "run", "main.go", "ls", "-bucket", bucketName)
		out, err := cmd.CombinedOutput()
		must.NoError(t, err)
		must.StrNotContains(t, string(out), "README.md")
	})
}
