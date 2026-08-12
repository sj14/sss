package e2e

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestCleanupMultipart(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("create mp in root", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "multipart", "create", "something")
		must.NoError(t, err)
	})

	t.Run("create mp in subdir/prefix", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "multipart", "create", "directory/something")
		must.NoError(t, err)
	})

	t.Run("list mp", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  directory/")
		must.StrContains(t, out, "something")
	})

	t.Run("cleanup dry-run", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "cleanup", "--force", "--all-multiparts", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list mp after dry-run", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrContains(t, out, "PREFIX  directory/")
		must.StrContains(t, out, "something")
	})

	t.Run("cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "cleanup", "--force", "--all-multiparts")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list mp after cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "multipart", "ls")
		must.NoError(t, err)
		must.StrNotContains(t, out, "PREFIX  directory/")
		must.StrNotContains(t, out, "something")
	})
}

func TestCleanupObjectVersions(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	bucketName := createBucket(t)

	t.Run("enable versioning", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "versioning", "put", "../samples/bucket-versioning.json")
		must.NoError(t, err)
	})

	// Two distinct keys, each with two versions, so a fix that mixes up
	// which key a version belongs to gets caught (wrong key/version-id
	// pairs get rejected by the backend).
	t.Run("upload object versions in root", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "something")
		must.NoError(t, err)
		_, err = run(t.Context(), "bucket", bucketName, "put", "../README.md", "something")
		must.NoError(t, err)
	})

	t.Run("upload object versions in subdir/prefix", func(t *testing.T) {
		_, err := run(t.Context(), "bucket", bucketName, "put", "../README.md", "directory/something")
		must.NoError(t, err)
		_, err = run(t.Context(), "bucket", bucketName, "put", "../README.md", "directory/something")
		must.NoError(t, err)
	})

	t.Run("list versions", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "versions", "--delimiter=''")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
		must.StrContains(t, out, "directory/something")
	})

	t.Run("cleanup dry-run", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "cleanup", "--force", "--all-object-versions", "--dry-run")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list versions after dry-run", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "versions", "--delimiter=''")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
		must.StrContains(t, out, "directory/something")
	})

	t.Run("cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "cleanup", "--force", "--all-object-versions")
		must.NoError(t, err)
		must.StrContains(t, out, "something")
	})

	t.Run("list versions after cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "versions", "--delimiter=''")
		must.NoError(t, err)
		must.StrNotContains(t, out, "something")
	})

	t.Run("list objects after cleanup", func(t *testing.T) {
		out, err := run(t.Context(), "bucket", bucketName, "ls", "--delimiter=''")
		must.NoError(t, err)
		must.StrNotContains(t, out, "something")
	})
}
