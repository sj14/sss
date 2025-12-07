.PHONY: test-deps
test-deps:
 	localstack start -e SERVICES=s3 -d
#   docker run -p 8333:8333 -e AWS_ACCESS_KEY_ID=your_access_key -e AWS_SECRET_ACCESS_KEY=your_secret_key chrislusf/seaweedfs server -s3

.PHONY: test-run
test-run:
	go test ./... -race
