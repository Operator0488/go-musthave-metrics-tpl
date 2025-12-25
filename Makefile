.PHONY: test1 test build clean
test1:
	./metricstest-darwin-arm64 -test.v -test.run=^TestIteration1$$ \
      -binary-path=cmd/server/server \
      -source-path=.

.PHONY: buildServer
buildServer:
	@go build -o cmd/server/server cmd/server/*.go

.PHONY: test
test:
	@go test