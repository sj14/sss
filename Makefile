.PHONY: test-deps
test-deps:
	localstack start -e SERVICES=s3 -d

.PHONY: test-run
test-run:
	SSS_CONFIG=./config_localstack.toml go test ./... -count=1 -race
