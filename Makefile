.PHONY: test-deps
test-deps:
	localstack start -e SERVICES=s3 -d
#   docker run -p 8333:8333 -e AWS_ACCESS_KEY_ID=your_access_key -e AWS_SECRET_ACCESS_KEY=your_secret_key chrislusf/seaweedfs server -s3
#   docker run --rm -e ROOT_ACCESS_KEY=testuser -e ROOT_SECRET_KEY=secret -e VGW_BACKEND=posix -e VGW_BACKEND_ARG=/tmp -p 7070:7070 versity/versitygw:latest

.PHONY: test-run
test-run:
	go test ./... -race -count=1 -timeout=1m
